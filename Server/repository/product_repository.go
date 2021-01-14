package repository

import (
	"encoding/json"

	"Taurant/driver"
	"Taurant/model"

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

func (b ProductRepositoryDGraph) FetchUIDs() (map[string]string, error) {
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
	res, err := driver.RunQuery(query)
	if err != nil {
		return nil, err
	}
	var productsFound []model.Product
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(res, &objmap)
	if err != nil {
		return nil, err
	}
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findAllProducts"], &productsFound)
	if err != nil {
		return nil, err
	}
	for _, productFound := range productsFound {
		iDUIDProducts[productFound.ID] = productFound.UID
	}
	return iDUIDProducts, nil
}

func (b ProductRepositoryDGraph) Create(product *model.Product) error {
	product.UID = "_:" + product.ID
	err := driver.RunMutation(product)
	if err != nil {
		return err
	}
	productFound, err := b.FindById(product.ID)
	if err != nil {
		return err
	}
	*(product) = *(productFound)
	return nil
}

func (b ProductRepositoryDGraph) Update(uid string, product *model.Product) error {
	product.UID = uid
	err := driver.RunMutation(product)
	if err != nil {
		return err
	}
	return nil
}

func (b ProductRepositoryDGraph) FindById(product_id string) (*model.Product, error) {
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
	res, err := driver.RunQuery(query)
	if err != nil {
		return nil, err
	}
	var productsFound []model.Product
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(res, &objmap)
	if err != nil {
		return nil, err
	}
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findProductById"], &productsFound)
	if err != nil {
		return nil, err
	}
	if len(productsFound) > 0 {
		return &productsFound[0], nil
	}
	return nil, nil
}
