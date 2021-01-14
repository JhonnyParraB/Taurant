package repository

import (
	"encoding/json"
	"fmt"

	"Taurant/driver"

	"Taurant/model"

	jsoniter "github.com/json-iterator/go"
)

//BuyerRepository _
type BuyerRepository interface {
	FetchBasicInformation() []model.Buyer
	Create(buyer *model.Buyer) string
	FindByID(buyerID string) *model.Buyer
	FindByIDWithTransactions(buyerID string) *model.Buyer
	Update(uid string, buyer *model.Buyer) string
	FindBuyersWithSameIP(buyerID string) []model.Buyer
}

//BuyerRepositoryDGraph _
type BuyerRepositoryDGraph struct {
}

//FetchBasicInformation _
func (b BuyerRepositoryDGraph) FetchBasicInformation(offset int64, itemsPerPage int64) ([]model.Buyer, error) {
	query := fmt.Sprintf(
		`
		{
			findAllBuyers(func: has(buyer_id), offset: %v, first: %v) {
				uid
				buyer_id
				buyer_name
				age
			}
		}	
	`, offset, itemsPerPage)
	res, err := driver.RunQuery(query)
	if err != nil {
		return nil, err
	}
	var buyersFound []model.Buyer
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(res, &objmap)
	if err != nil {
		return nil, err
	}
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findAllBuyers"], &buyersFound)
	if err != nil {
		return nil, err
	}
	return buyersFound, nil
}

//Create _
func (b BuyerRepositoryDGraph) Create(buyer *model.Buyer) error {
	buyer.UID = "_:" + buyer.ID
	err := driver.RunMutation(buyer)
	if err != nil {
		return err
	}
	buyerFound, err := b.FindByID(buyer.ID)
	if err != nil {
		return err
	}
	*(buyer) = *(buyerFound)
	return nil
}

//Update _
func (b BuyerRepositoryDGraph) Update(uid string, buyer *model.Buyer) error {
	buyer.UID = uid
	err := driver.RunMutation(buyer)
	if err != nil {
		return err
	}
	return nil
}

//FindByID _
func (b BuyerRepositoryDGraph) FindByID(buyerID string) (*model.Buyer, error) {
	query :=
		`
		{
			findBuyerById(func: eq(buyer_id, "` + buyerID + `"), first: 1) {
				uid
				buyer_id
				buyer_name
				age
			}
		}
	`
	res, err := driver.RunQuery(query)
	if err != nil {
		return nil, err
	}
	var buyersFound []model.Buyer
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(res, &objmap)
	if err != nil {
		return nil, err
	}
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findBuyerById"], &buyersFound)
	if err != nil {
		return nil, err
	}
	if len(buyersFound) > 0 {
		return &buyersFound[0], nil
	}
	return nil, nil
}

//FindByIDWithTransactions _
func (b BuyerRepositoryDGraph) FindByIDWithTransactions(buyerID string) (*model.Buyer, error) {
	query :=
		`
		{
			findBuyerById(func: eq(buyer_id, "` + buyerID + `"), first: 1) {
				uid
				buyer_id
				buyer_name
				age
				perform: ~is_made_by{
					transaction_id
					transaction_date
					location{
						ip
					}
					device
					include {
						trade{
							product_id
							product_name
							price
						}
						quantity
					}
				}
			}
		}	
	`
	res, err := driver.RunQuery(query)
	if err != nil {
		return nil, err
	}
	var buyersFound []model.Buyer
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(res, &objmap)
	if err != nil {
		return nil, err
	}
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findBuyerById"], &buyersFound)
	if err != nil {
		return nil, err
	}
	if len(buyersFound) > 0 {
		return &buyersFound[0], nil
	}
	return nil, nil
}

//FindBuyersWithSameIP _
func (b BuyerRepositoryDGraph) FindBuyersWithSameIP(buyerID string) ([]model.Buyer, error) {
	query :=
		`
		{
			var(func:eq(buyer_id, "` + buyerID + `")){
			  transactions: ~is_made_by{
				  m as location{
				   ip
				   has_transactions:~location @groupby(is_made_by){
					  a as count(uid)
					}
				}
				}
			}
			
			BuyersWithSameIP(func:uid(a)) @cascade @normalize @filter(not(eq(buyer_id, "` + buyerID + `"))){
			  uid: uid
			  buyer_id: buyer_id 
			  buyer_name: buyer_name
			  age: age
			  ~is_made_by {
				location @filter(uid(m)){
				  shared_ip: ip
				}
			  }
			} 
		}	
	`
	res, err := driver.RunQuery(query)
	if err != nil {
		return nil, err
	}
	var buyersWithSameIP []model.Buyer
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(res, &objmap)
	if err != nil {
		return nil, err
	}
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["BuyersWithSameIP"], &buyersWithSameIP)
	if err != nil {
		return nil, err
	}
	return buyersWithSameIP, nil
}

//FindRecommendedProducts _
func (b BuyerRepositoryDGraph) FindRecommendedProducts(buyerID string) ([]model.Product, error) {
	query :=
		`
		{
			var(func:has(quantity)) @groupby(trade) {
				times_buyed as sum(quantity)
			}
	
			var(func:eq(buyer_id, "` + buyerID + `")){
				~is_made_by{
					include{
						products_previously_buyed as trade
					}
				}
			}
	
			FindRecommendedProducts(func: uid(times_buyed), orderdesc: val(times_buyed), first:5) @filter(not(uid(products_previously_buyed))){
				uid
				product_id 
				product_name
				price
			}
  		}	
	`
	res, err := driver.RunQuery(query)
	if err != nil {
		return nil, err
	}
	var recommendedProducts []model.Product
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(res, &objmap)
	if err != nil {
		return nil, err
	}
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["FindRecommendedProducts"], &recommendedProducts)
	if err != nil {
		return nil, err
	}
	return recommendedProducts, nil
}
