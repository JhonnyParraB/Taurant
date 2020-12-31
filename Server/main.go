package main

import (
	"io/ioutil"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"./model"
	"./repository"
)

func main() {
	var endpointCaseJSON = jsoniter.Config{TagKey: "endpoint"}.Froze()
	var buyer_repository repository.BuyerRepositoryDGraph

	response, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers")
	handleError(err)
	data, _ := ioutil.ReadAll(response.Body)
	var buyers []model.Buyer
	err = endpointCaseJSON.Unmarshal(data, &buyers)
	for _, buyer := range buyers {
		buyer_found := buyer_repository.FindById(buyer.ID)
		if buyer_found == nil {
			buyer_repository.Create(&buyer)
		} else {
			buyer_repository.Update(buyer_found.UID, &buyer)
		}
	}
}

func handleError(err interface{}) {
	if err != nil {
		log.Fatal(err)
	}
}
