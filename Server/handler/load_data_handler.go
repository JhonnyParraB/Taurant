package handler

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"../driver"
	"../model"
	"../repository"
	"github.com/dgraph-io/dgo"
	"github.com/go-chi/chi"
	jsoniter "github.com/json-iterator/go"
)

type LoadDayDataHandler struct {
	buyerRepository       repository.BuyerRepositoryDGraph
	productRepository     repository.ProductRepositoryDGraph
	transactionRepository repository.TransactionRepositoryDGraph
	locationRepository    repository.LocationRepositoryDGraph
	jobChan               chan loadDataJob
	iDUIDBuyers           map[string]string
	iDUIDProducts         map[string]string
	iDUIDLocations        map[string]string
}

type loadDataJob struct {
	date int
}

func (l *LoadDayDataHandler) Init() {
	l.jobChan = make(chan loadDataJob, 100)
	go l.worker(l.jobChan)
	l.iDUIDBuyers = make(map[string]string)
	l.iDUIDProducts = make(map[string]string)
	l.iDUIDLocations = make(map[string]string)
}

func (l *LoadDayDataHandler) worker(jobChan <-chan loadDataJob) {
	for job := range jobChan {
		l.processLoadDataJob(job)
	}
}

func (l *LoadDayDataHandler) processLoadDataJob(job loadDataJob) {
	txn := driver.CreateTransaction()
	defer txn.Discard(context.Background())
	l.loadBuyers(txn, job.date)
	l.loadProducts(txn, job.date)
	l.loadTransactions(txn, job.date)
	txn.Commit(context.Background())
}

func (l *LoadDayDataHandler) LoadDayData(w http.ResponseWriter, r *http.Request) {
	dateParam := chi.URLParam(r, "date")
	dateUnixFormat, err := strconv.Atoi(dateParam)
	handleError(err)
	a := true
	job := loadDataJob{date: dateUnixFormat}
	select {
	case l.jobChan <- job:
		a = true
	default:
		a = false
	}
	print(a)
}

func (l *LoadDayDataHandler) loadBuyers(txn *dgo.Txn, date int) {
	var endpointCaseJSON = jsoniter.Config{TagKey: "endpoint"}.Froze()
	response, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers")
	handleError(err)
	data, _ := ioutil.ReadAll(response.Body)
	var buyers []model.Buyer
	err = endpointCaseJSON.Unmarshal(data, &buyers)
	for _, buyer := range buyers {
		buyerFound := l.buyerRepository.FindById(buyer.ID)
		if buyerFound == nil {
			l.buyerRepository.AddCreateToTransaction(txn, &buyer)
		} else {
			l.buyerRepository.AddUpdateToTransaction(txn, buyerFound.UID, &buyer)
		}
		l.iDUIDBuyers[buyer.ID] = buyer.UID
	}
}

func (l *LoadDayDataHandler) loadProducts(txn *dgo.Txn, date int) {
	response, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products")
	handleError(err)
	data, _ := ioutil.ReadAll(response.Body)
	bytesReader := bytes.NewReader(data)
	csvReader := csv.NewReader(bytesReader)
	csvReader.Comma = '\''
	for {
		recordFields, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		price, err := strconv.Atoi(recordFields[2])
		handleError(err)
		product := model.Product{
			ID:    recordFields[0],
			Name:  recordFields[1],
			Price: price,
		}
		productFound := l.productRepository.FindById(product.ID)
		if productFound == nil {
			l.productRepository.AddCreateToTransaction(txn, &product)
		} else {
			l.productRepository.AddUpdateToTransaction(txn, productFound.UID, &product)
		}
		l.iDUIDProducts[product.ID] = product.UID
	}
}

func (l *LoadDayDataHandler) loadTransactions(txn *dgo.Txn, date int) {
	response, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions")
	data, _ := ioutil.ReadAll(response.Body)
	bytesReader := bytes.NewReader(data)
	scanner := bufio.NewScanner(bytesReader)
	scanner.Split(splitAt("\000\000"))
	a := 0
	for scanner.Scan() {
		record := scanner.Text()
		recordFields := strings.Split(record, "\000")

		a++
		if a%500 == 0 {
			fmt.Printf("%d   --   "+record+"\n", a)
		}

		//set buyer
		buyerID := recordFields[1]
		var buyer model.Buyer
		if uid, exists := l.iDUIDBuyers[buyerID]; exists {
			buyer = model.Buyer{

				UID: uid,
			}
		}

		//set products
		recordFields[4] = strings.Replace(recordFields[4], "(", "", -1)
		recordFields[4] = strings.Replace(recordFields[4], ")", "", -1)
		productsID := strings.Split(recordFields[4], ",")
		var products []model.Product
		var productOrders []model.ProductOrder
		for _, productID := range productsID {
			if uid, exists := l.iDUIDProducts[productID]; exists {
				products = append(products, model.Product{
					UID: uid,
				})
			}
		}
		dupMap := dupCount(products)
		for uid, quantity := range dupMap {
			productOrders = append(productOrders, model.ProductOrder{
				Product: &model.Product{
					UID: uid,
				},
				Quantity: quantity,
			})
		}

		//set location
		IP := recordFields[2]
		location := model.Location{
			IP: IP,
		}
		locationFound := l.locationRepository.FindByIP(IP)
		if locationFound == nil {
			l.locationRepository.AddCreateToTransaction(txn, &location)
		} else {
			location.UID = locationFound.UID
		}

		transaction := model.Transaction{
			ID:            recordFields[0],
			Buyer:         &buyer,
			Location:      &location,
			Device:        recordFields[3],
			ProductOrders: &productOrders,
		}

		transactionFound := l.transactionRepository.FindById(transaction.ID)
		if transactionFound == nil {
			l.transactionRepository.AddCreateToTransaction(txn, &transaction)
		} else {
			l.transactionRepository.AddUpdateToTransaction(txn, transactionFound.UID, &transaction)
		}
	}
	handleError(err)
}

func splitAt(substring string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	searchBytes := []byte(substring)
	searchLen := len(searchBytes)
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		dataLen := len(data)

		// Return nothing if at end of file and no data passed
		if atEOF && dataLen == 0 {
			return 0, nil, nil
		}

		// Find next separator and return token
		if i := bytes.Index(data, searchBytes); i >= 0 {
			return i + searchLen, data[0:i], nil
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return dataLen, data, nil
		}

		// Request more data.
		return 0, nil, nil
	}
}

func dupCount(list []model.Product) map[string]int {

	duplicateFrequency := make(map[string]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		_, exist := duplicateFrequency[item.UID]

		if exist {
			duplicateFrequency[item.UID]++ // increase counter by 1 if already in the map
		} else {
			duplicateFrequency[item.UID] = 1 // else start counting from 1
		}
	}
	return duplicateFrequency
}

func handleError(err interface{}) {
	if err != nil {
		log.Fatal(err)
	}
}
