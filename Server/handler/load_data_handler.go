package handler

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"time"

	"Taurant/model"
	"Taurant/repository"

	"github.com/go-chi/chi"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/ksuid"
)

const errorTag string = "ERROR:"
const errorEmail = "ERROR"
const succesfulEmail = "SUCCESFUL"
const buyersExternalEndpoint string = "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers"
const productsExternalEndpoint string = "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products"
const transactionsExternalEndpoint string = "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions"

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//LoadDayDataHandler is used at router
type LoadDayDataHandler struct {
	buyerRepository       repository.BuyerRepositoryDGraph
	productRepository     repository.ProductRepositoryDGraph
	transactionRepository repository.TransactionRepositoryDGraph
	locationRepository    repository.LocationRepositoryDGraph
	laodDataJobRepository repository.LoadDataJobRepositoryDGraph
	jobChan               chan model.LoadDataJob
	iDUIDBuyers           map[string]string
	iDUIDProducts         map[string]string
	iDUIDLocations        map[string]string
}

//NewLoadDayDataHandler is used at router
func NewLoadDayDataHandler() *LoadDayDataHandler {
	var l LoadDayDataHandler
	l.jobChan = make(chan model.LoadDataJob, 100)
	jobsPending, err := l.laodDataJobRepository.FetchBasicInformation()
	if err != nil {
		log.Fatal(errorTag, "Error when trying to bring pending jobs from the database", err)
	}
	for _, jobPending := range jobsPending {
		l.jobChan <- jobPending
	}
	go l.worker(l.jobChan)
	l.iDUIDBuyers = make(map[string]string)
	l.iDUIDProducts = make(map[string]string)
	l.iDUIDLocations = make(map[string]string)
	return &l
}

func (l *LoadDayDataHandler) worker(jobChan <-chan model.LoadDataJob) {
	for job := range jobChan {
		l.processLoadDataJob(job)
	}
}

func (l *LoadDayDataHandler) processLoadDataJob(job model.LoadDataJob) {
	log.Println("STARTING THE DATA LOAD JOB ... ", job.ID)
	thereIsEmail := job.Email != ""
	err := l.loadBuyers(job.Date)
	if err != nil {
		handleLoadDataJobError(err, thereIsEmail, job)
		return
	}
	err = l.loadProducts(job.Date)
	if err != nil {
		handleLoadDataJobError(err, thereIsEmail, job)
		return
	}
	err = l.loadTransactions(job.Date)
	if err != nil {
		handleLoadDataJobError(err, thereIsEmail, job)
		return
	}
	if thereIsEmail {
		sendLoadDataResultEmail(job, getSuccessEmailMessage(job), succesfulEmail)
	}
	err = l.laodDataJobRepository.Delete(&job)
	if err != nil {
		log.Println(errorTag, "Job completed but there was an error removing it from the database", job.ID, err)
		return
	}
}

func handleLoadDataJobError(err error, thereIsEmail bool, job model.LoadDataJob) {
	log.Println(errorTag, job.ID, err)
	if thereIsEmail {
		sendLoadDataResultEmail(job, getErrorEmailMessage(job), errorEmail)
	}
}

func getSuccessEmailMessage(job model.LoadDataJob) []byte {
	tm := time.Unix(job.Date, 0)
	msg := []byte(fmt.Sprintf("To: %s\r\n", job.Email) +
		fmt.Sprintf("Subject: Load Data Request %s succesfully completed!\r\n", job.ID) +
		"\r\n" +
		fmt.Sprintf("Your load data request for %v (%s) time was succesfully completed.\r\n -Taurant\r\n", job.Date, tm.Format(time.UnixDate)))
	return msg
}

func getErrorEmailMessage(job model.LoadDataJob) []byte {
	tm := time.Unix(job.Date, 0)
	msg := []byte(fmt.Sprintf("To: %s\r\n", job.Email) +
		fmt.Sprintf("Subject: Error While Processing Load Data Request %s\r\n", job.ID) +
		"\r\n" +
		fmt.Sprintf("There was an error while processing your load data request for %v (%s). Please contact Taurant's Administrator.\r\n -Taurant\r\n", job.Date, tm.Format(time.UnixDate)))
	return msg
}

func sendLoadDataResultEmail(job model.LoadDataJob, msg []byte, emailType string) {
	from := "ttaurant@gmail.com"
	password := "taurant123456"

	to := []string{
		job.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		log.Println(errorTag, err)
	}
	log.Println("EmailSended: ", job.Email, job.ID, emailType)
}

//LoadDayData is used at router
func (l *LoadDayDataHandler) LoadDayData(w http.ResponseWriter, r *http.Request) {
	dateParam := chi.URLParam(r, "date")
	email := r.URL.Query().Get("email")
	if email != "" {
		if !isEmailValid(email) {
			log.Println(errorTag, r.Method, r.URL.String(), "Wrong email format")
			respondwithJSON(w, http.StatusBadRequest,
				errorMessage{
					Code:    wrongEmailParamFormatCode,
					Details: errorsDetails[wrongEmailParamFormatCode],
				})
			return
		}
	}
	dateUnixFormat, err := strconv.ParseInt(dateParam, 10, 64)
	if err != nil {
		log.Println(errorTag, r.Method, r.URL.String(), err)
		respondwithJSON(w, http.StatusBadRequest,
			errorMessage{
				Code:    wrongDateParamFormatCode,
				Details: errorsDetails[wrongDateParamFormatCode],
			})
		return
	}
	ksuid := ksuid.New().String()
	job := model.LoadDataJob{
		UID:   "_:" + ksuid,
		Date:  dateUnixFormat,
		ID:    ksuid,
		Email: email,
	}
	err = l.laodDataJobRepository.Create(&job)
	if err != nil {
		log.Println(errorTag, r.Method, r.URL.String(), "An error occurred when trying to save the job in the database", err)
		respondwithJSON(w, http.StatusInternalServerError,
			errorMessage{
				Code:    internalServerErrorCode,
				Details: errorsDetails[internalServerErrorCode],
			})
		return
	}
	select {
	case l.jobChan <- job:
		respondwithJSON(w, http.StatusOK,
			responseMessage{
				Details: fmt.Sprintf("Request for load data for %v time has been received", job.Date),
				ID:      job.ID,
			})
	default:
		log.Println(errorTag, r.Method, r.URL.String(), "The job queue for load data is full")
		respondwithJSON(w, http.StatusInternalServerError,
			errorMessage{
				Code:    queueForLoadDataJobFullCode,
				Details: errorsDetails[queueForLoadDataJobFullCode],
			})
		return
	}
}

func (l *LoadDayDataHandler) loadBuyers(date int64) error {
	var endpointCaseJSON = jsoniter.Config{TagKey: "endpoint"}.Froze()
	response, err := http.Get(buyersExternalEndpoint + fmt.Sprintf("?date=%v", date))
	if err != nil {
		return err
	}
	data, _ := ioutil.ReadAll(response.Body)
	var buyers []model.Buyer
	err = endpointCaseJSON.Unmarshal(data, &buyers)
	for _, buyer := range buyers {
		buyerFound, err := l.buyerRepository.FindById(buyer.ID)
		if err != nil {
			return err
		}
		if buyerFound == nil {
			err = l.buyerRepository.Create(&buyer)
			if err != nil {
				return err
			}
		} else {
			err = l.buyerRepository.Update(buyerFound.UID, &buyer)
			if err != nil {
				return err
			}
		}
		l.iDUIDBuyers[buyer.ID] = buyer.UID
	}
	return nil
}

func (l *LoadDayDataHandler) loadProducts(date int64) error {
	response, err := http.Get(productsExternalEndpoint + fmt.Sprintf("?date=%v", date))
	if err != nil {
		return err
	}
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
		if err != nil {
			return err
		}
		product := model.Product{
			ID:    recordFields[0],
			Name:  recordFields[1],
			Price: price,
		}
		productFound, err := l.productRepository.FindById(product.ID)
		if err != nil {
			return err
		}
		if productFound == nil {
			err = l.productRepository.Create(&product)
			if err != nil {
				return err
			}
		} else {
			err = l.productRepository.Update(productFound.UID, &product)
			if err != nil {
				return err
			}
		}
		l.iDUIDProducts[product.ID] = product.UID
	}
	return nil
}

func (l *LoadDayDataHandler) loadTransactions(date int64) error {
	response, err := http.Get(transactionsExternalEndpoint + fmt.Sprintf("?date=%v", date))
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	bytesReader := bytes.NewReader(data)
	scanner := bufio.NewScanner(bytesReader)
	scanner.Split(splitAt("\000\000"))
	for scanner.Scan() {
		record := scanner.Text()
		recordFields := strings.Split(record, "\000")

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
				UID: "_:" + ksuid.New().String(),
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
		locationFound, err := l.locationRepository.FindByIP(IP)
		if err != nil {
			return err
		}
		if locationFound == nil {
			err = l.locationRepository.Create(&location)
			if err != nil {
				return err
			}
		} else {
			location.UID = locationFound.UID
		}

		transaction := model.Transaction{
			ID:            recordFields[0],
			Buyer:         &buyer,
			Date:          date,
			Location:      &location,
			Device:        recordFields[3],
			ProductOrders: &productOrders,
		}

		transactionFound, err := l.transactionRepository.FindById(transaction.ID)
		if err != nil {
			return err
		}
		if transactionFound == nil {
			err = l.transactionRepository.Create(&transaction)
			if err != nil {
				return err
			}
		} /*else {
			for _, productOrder := range *transactionFound.ProductOrders {
				err = l.transactionRepository.DeleteProductOrder(&productOrder)
				if err != nil {
					return err
				}
			}
			err = l.transactionRepository.Update(transactionFound.UID, &transaction)
			if err != nil {
				return err
			}

		}*/
		//IT WAS NOT POSSIBLE TO DELETE OLD PRODUCT_ORDERS
	}
	return nil
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
