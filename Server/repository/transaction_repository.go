package repository

import (
	"encoding/json"

	"Taurant/driver"
	"Taurant/model"

	jsoniter "github.com/json-iterator/go"
)

//TransactionRepository _
type TransactionRepository interface {
	Create(transaction *model.Transaction) string
	FindByID(transactionID string) *model.Transaction
	Update(uid string, transaction *model.Transaction) string
	DeleteProductOrders(productOrders []*model.ProductOrder)
}

//TransactionRepositoryDGraph _
type TransactionRepositoryDGraph struct {
}

//Create _
func (b TransactionRepositoryDGraph) Create(transaction *model.Transaction) error {
	transaction.UID = "_:" + transaction.ID
	err := driver.RunMutation(transaction)
	if err != nil {
		return err
	}
	transactionFound, err := b.FindByID(transaction.ID)
	if err != nil {
		return err
	}
	*(transaction) = *(transactionFound)
	return nil
}

//Update _
func (b TransactionRepositoryDGraph) Update(uid string, transaction *model.Transaction) error {
	transaction.UID = uid
	err := driver.RunMutation(transaction)
	if err != nil {
		return err
	}
	return nil
}

//FindByID _
func (b TransactionRepositoryDGraph) FindByID(transactionID string) (*model.Transaction, error) {
	query :=
		`
		{
			findTransactionById(func: eq(transaction_id, "` + transactionID + `"), first: 1) {
				uid
				transaction_id
				transaction_date
				is_made_by {
					uid
					id
					name
				}
				device
				ip
				include{
					uid
				}
			}
		}	
	`
	res, err := driver.RunQuery(query)
	if err != nil {
		return nil, err
	}
	var transactionsFound []model.Transaction
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(res, &objmap)
	if err != nil {
		return nil, err
	}
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findTransactionById"], &transactionsFound)
	if err != nil {
		return nil, err
	}
	if len(transactionsFound) > 0 {
		return &transactionsFound[0], nil
	}
	return nil, nil
}

//DeleteProductOrder _
func (b TransactionRepositoryDGraph) DeleteProductOrder(productOrder *model.ProductOrder) error {
	err := driver.RunMutationForDelete(map[string]string{"uid": productOrder.UID})
	if err != nil {
		return err
	}
	return nil
}
