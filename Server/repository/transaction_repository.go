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

func (b TransactionRepositoryDGraph) Create(transaction *model.Transaction) error {
	transaction.UID = "_:" + transaction.ID
	err := driver.RunMutation(transaction)
	if err != nil {
		return err
	}
	transaction, err = b.FindById(transaction.ID)
	if err != nil {
		return err
	}
	return nil
}

func (b TransactionRepositoryDGraph) Update(uid string, transaction *model.Transaction) error {
	transaction.UID = uid
	err := driver.RunMutation(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (b TransactionRepositoryDGraph) FindById(transaction_id string) (*model.Transaction, error) {
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
