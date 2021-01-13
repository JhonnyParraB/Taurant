package repository

import (
	"encoding/json"

	"../driver"
	"../model"
	"github.com/dgraph-io/dgo"
	jsoniter "github.com/json-iterator/go"
)

type LocationRepository interface {
	Create(location *model.Location) string
	FindByIP(ip string) *model.Location
	Update(uid string, location *model.Location) string
}

type LocationRepositoryDGraph struct {
}

func (b LocationRepositoryDGraph) Create(location *model.Location) {
	location.UID = "_:" + location.IP
	driver.RunMutation(location)
	*location = *(b.FindByIP(location.IP))
}

func (b LocationRepositoryDGraph) Update(uid string, location *model.Location) {
	location.UID = uid
	driver.RunMutation(location)
}

func (b LocationRepositoryDGraph) AddCreateToTransaction(txn *dgo.Txn, location *model.Location) {
	location.UID = "_:" + location.IP
	driver.AddMutationToTransaction(txn, location)
}

func (b LocationRepositoryDGraph) AddUpdateToTransaction(txn *dgo.Txn, uid string, location *model.Location) {
	location.UID = uid
	driver.AddMutationToTransaction(txn, location)
}

func (b LocationRepositoryDGraph) FindByIP(ip string) *model.Location {
	query :=
		`
		{
			findLocationByIP(func: eq(ip, "` + ip + `"), first: 1) {
				uid
				ip
			}
		}	
	`
	res := driver.RunQuery(query)
	var locationsFound []model.Location
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(res, &objmap)
	handleError(err)
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findLocationByIP"], &locationsFound)
	handleError(err)
	if len(locationsFound) > 0 {
		return &locationsFound[0]
	}
	return nil
}
