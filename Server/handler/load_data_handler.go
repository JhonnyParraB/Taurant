package handler

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../model"
	"../repository"
	jsoniter "github.com/json-iterator/go"
)

var iDUIDBuyers map[string]string
var iDUIDProducts map[string]string
var iDUIDLocations map[string]string

func init() {
	iDUIDBuyers = make(map[string]string)
	iDUIDProducts = make(map[string]string)
	iDUIDLocations = make(map[string]string)
}

type LoadDayDataHandler struct {
	buyerRepository       repository.BuyerRepositoryDGraph
	productRepository     repository.ProductRepositoryDGraph
	transactionRepository repository.TransactionRepositoryDGraph
	locationRepository    repository.LocationRepositoryDGraph
}

func (l *LoadDayDataHandler) LoadDayData(w http.ResponseWriter, r *http.Request) {
	//date, err := strconv.Atoi(r.URL.Query().Get("date"))
	date := int32(time.Now().Unix())
	l.loadBuyers(date)
	l.loadProducts(date)
	l.loadTransactions(date)
}

func (l *LoadDayDataHandler) loadBuyers(date int32) {
	var endpointCaseJSON = jsoniter.Config{TagKey: "endpoint"}.Froze()
	response, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers")
	handleError(err)
	data, _ := ioutil.ReadAll(response.Body)
	var buyers []model.Buyer
	err = endpointCaseJSON.Unmarshal(data, &buyers)
	for _, buyer := range buyers {
		buyerFound := l.buyerRepository.FindById(buyer.ID)
		if buyerFound == nil {
			l.buyerRepository.Create(&buyer)
		} else {
			l.buyerRepository.Update(buyerFound.UID, &buyer)
		}
		iDUIDBuyers[buyer.ID] = buyer.UID
	}
}

func (l *LoadDayDataHandler) loadProducts(date int32) {
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
			l.productRepository.Create(&product)
		} else {
			l.productRepository.Update(productFound.UID, &product)
		}
		iDUIDProducts[product.ID] = product.UID
	}
}

func (l *LoadDayDataHandler) loadTransactions(date int32) {
	response, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions")
	data, _ := ioutil.ReadAll(response.Body)
	bytesReader := bytes.NewReader(data)
	scanner := bufio.NewScanner(bytesReader)
	scanner.Split(splitAt("\000\000"))
	for scanner.Scan() {
		record := scanner.Text()
		recordFields := strings.Split(record, "\000")

		//set buyer
		buyerID := recordFields[1]
		var buyer model.Buyer
		if uid, exists := iDUIDBuyers[buyerID]; exists {
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
			if uid, exists := iDUIDProducts[productID]; exists {
				products = append(products, model.Product{
					UID: uid,
				})
			}
		}
		dupMap := dup_count(products)
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
			l.locationRepository.Create(&location)
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
			l.transactionRepository.Create(&transaction)
		} else {
			l.transactionRepository.Update(transactionFound.UID, &transaction)
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

func dup_count(list []model.Product) map[string]int {

	duplicate_frequency := make(map[string]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		_, exist := duplicate_frequency[item.UID]

		if exist {
			duplicate_frequency[item.UID] += 1 // increase counter by 1 if already in the map
		} else {
			duplicate_frequency[item.UID] = 1 // else start counting from 1
		}
	}
	return duplicate_frequency
}

func handleError(err interface{}) {
	if err != nil {
		log.Fatal(err)
	}
}
