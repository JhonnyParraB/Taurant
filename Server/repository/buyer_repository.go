package repository

import (
	"encoding/json"
	"log"

	"../driver"
	"../model"
	jsoniter "github.com/json-iterator/go"
)

type BuyerRepository interface {
	FetchUIDs() map[string]string
	Create(buyer *model.Buyer) string
	FindById(buyer_id string) *model.Buyer
	Update(uid string, buyer *model.Buyer) string
}

type BuyerRepositoryDGraph struct {
}

func (b BuyerRepositoryDGraph) FetchUIDs() map[string]string {
	iDUIDBuyers := make(map[string]string)
	query :=
		`
		{
			findAllBuyers(func: has(buyer_id)) {
				uid
				buyer_id
			}
		}	
	`
	res := driver.RunQuery(query)
	var buyersFound []model.Buyer
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(res, &objmap)
	handleError(err)
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findAllBuyers"], &buyersFound)
	for _, buyerFound := range buyersFound {
		iDUIDBuyers[buyerFound.ID] = buyerFound.UID
	}
	return iDUIDBuyers
}

func (b BuyerRepositoryDGraph) Create(buyer *model.Buyer) {
	buyer.UID = "_:" + buyer.ID
	driver.RunMutation(buyer)
}

func (b BuyerRepositoryDGraph) Update(uid string, buyer *model.Buyer) string {
	buyer.UID = uid
	driver.RunMutation(buyer)
	return buyer.ID
}

func (b BuyerRepositoryDGraph) FindById(buyer_id string) *model.Buyer {
	query :=
		`
		{
			findBuyerById(func: eq(buyer_id, "` + buyer_id + `"), first: 1) {
				uid
				buyer_id
				buyer_name
				age
				transactions
			}
		}	
	`
	res := driver.RunQuery(query)
	var buyersFound []model.Buyer
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(res, &objmap)
	handleError(err)
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findBuyerById"], &buyersFound)
	handleError(err)
	if len(buyersFound) > 0 {
		return &buyersFound[0]
	}
	return nil
}

func handleError(err interface{}) {
	if err != nil {
		log.Fatal(err)
	}
}
