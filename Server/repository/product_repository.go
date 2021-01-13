package repository

import (
	"encoding/json"

	"../driver"
	"../model"
	"github.com/dgraph-io/dgo"
	jsoniter "github.com/json-iterator/go"
)

type ProductRepository interface {
	FetchUIDs() map[string]string
	Create(product *model.Product) string
	FindById(product_id string) *model.Product
	Update(uid string, product *model.Product) string
}

type ProductRepositoryDGraph struct {
}

func (b ProductRepositoryDGraph) FetchUIDs() map[string]string {
	iDUIDProducts := make(map[string]string)
	query :=
		`
		{
			findAllProducts(func: has(product_id)) {
				uid
				product_id
			}
		}	
	`
	res := driver.RunQuery(query)
	var productsFound []model.Product
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(res, &objmap)
	handleError(err)
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findAllProducts"], &productsFound)
	for _, productFound := range productsFound {
		iDUIDProducts[productFound.ID] = productFound.UID
	}
	return iDUIDProducts
}

func (b ProductRepositoryDGraph) Create(product *model.Product) {
	product.UID = "_:" + product.ID
	driver.RunMutation(product)
	*product = *(b.FindById(product.ID))
}

func (b ProductRepositoryDGraph) Update(uid string, product *model.Product) {
	product.UID = uid
	driver.RunMutation(product)
}

func (b ProductRepositoryDGraph) AddCreateToTransaction(txn *dgo.Txn, product *model.Product) {
	product.UID = "_:" + product.ID
	driver.AddMutationToTransaction(txn, product)
}

func (b ProductRepositoryDGraph) AddUpdateToTransaction(txn *dgo.Txn, uid string, product *model.Product) {
	product.UID = uid
	driver.AddMutationToTransaction(txn, product)
}

func (b ProductRepositoryDGraph) FindById(product_id string) *model.Product {
	query :=
		`
		{
			findProductById(func: eq(product_id, "` + product_id + `"), first: 1) {
				uid
				product_id
				product_name
				price
				has_transactions
			}
		}	
	`
	res := driver.RunQuery(query)
	var productsFound []model.Product
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(res, &objmap)
	handleError(err)
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findProductById"], &productsFound)
	handleError(err)
	if len(productsFound) > 0 {
		return &productsFound[0]
	}
	return nil
}
