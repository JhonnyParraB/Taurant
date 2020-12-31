package repository

import (
	"encoding/json"
	"log"

	"../driver"
	"../model"
)

type BuyerRepository interface {
	Create(buyer *model.Buyer) string
	FindById(buyer_id string) *model.Buyer
	Update(uid string, buyer *model.Buyer) string
}

type BuyerRepositoryDGraph struct {
}

func (b BuyerRepositoryDGraph) Create(buyer *model.Buyer) string {
	buyer.UID = "_:" + buyer.ID
	driver.RunMutation(buyer)
	return buyer.ID
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
	err = json.Unmarshal(objmap["findBuyerById"], &buyersFound)
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
