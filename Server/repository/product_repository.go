package repository

import (
	"encoding/json"

	"../driver"
	"../model"
)

type ProductRepository interface {
	Create(product *model.Product) string
	FindById(product_id string) *model.Product
	Update(uid string, product *model.Product) string
}

type ProductRepositoryDGraph struct {
}

func (b ProductRepositoryDGraph) Create(product *model.Product) string {
	product.UID = "_:" + product.ID
	driver.RunMutation(product)
	return product.ID
}

func (b ProductRepositoryDGraph) Update(uid string, product *model.Product) string {
	product.UID = uid
	driver.RunMutation(product)
	return product.ID
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
	err = json.Unmarshal(objmap["findProductById"], &productsFound)
	handleError(err)
	if len(productsFound) > 0 {
		return &productsFound[0]
	}
	return nil
}
