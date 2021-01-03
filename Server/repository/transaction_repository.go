package repository

import (
	"encoding/json"

	"../driver"
	"../model"
	jsoniter "github.com/json-iterator/go"
)

type TransactionRepository interface {
	Create(transaction *model.Transaction) string
	FindById(transaction_id string) *model.Transaction
	Update(uid string, transaction *model.Transaction) string
}

type TransactionRepositoryDGraph struct {
}

func (b TransactionRepositoryDGraph) Create(transaction *model.Transaction) {
	transaction.UID = "_:" + transaction.ID
	driver.RunMutation(transaction)
	*transaction = *(b.FindById(transaction.ID))
}

func (b TransactionRepositoryDGraph) Update(uid string, transaction *model.Transaction) {
	transaction.UID = uid
	driver.RunMutation(transaction)
}

func (b TransactionRepositoryDGraph) FindById(transaction_id string) *model.Transaction {
	query :=
		`
		{
			findTransactionById(func: eq(transaction_id, "` + transaction_id + `"), first: 1) {
				uid
				transaction_id
				is_made_by
				device
				ip
				trade
			}
		}	
	`
	res := driver.RunQuery(query)
	var transactionsFound []model.Transaction
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(res, &objmap)
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findTransactionById"], &transactionsFound)
	handleError(err)
	if len(transactionsFound) > 0 {
		return &transactionsFound[0]
	}
	return nil
}
