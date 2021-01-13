package repository

import (
	"encoding/json"

	"../driver"
	"../model"
	jsoniter "github.com/json-iterator/go"
)

type LocationRepository interface {
	Create(location *model.Location) string
	FindByIP(ip string) *model.Location
	Update(uid string, location *model.Location) string
}

type LocationRepositoryDGraph struct {
}

func (b LocationRepositoryDGraph) Create(location *model.Location) error {
	location.UID = "_:" + location.IP
	err := driver.RunMutation(location)
	if err != nil {
		return err
	}
	location, err = b.FindByIP(location.IP)
	if err != nil {
		return err
	}
	return nil
}

func (b LocationRepositoryDGraph) Update(uid string, location *model.Location) error {
	location.UID = uid
	err := driver.RunMutation(location)
	if err != nil {
		return err
	}
	return nil
}

func (b LocationRepositoryDGraph) FindByIP(ip string) (*model.Location, error) {
	query :=
		`
		{
			findLocationByIP(func: eq(ip, "` + ip + `"), first: 1) {
				uid
				ip
			}
		}	
	`
	res, err := driver.RunQuery(query)
	if err != nil {
		return nil, err
	}
	var locationsFound []model.Location
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(res, &objmap)
	if err != nil {
		return nil, err
	}
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	err = predicateCaseJSON.Unmarshal(objmap["findLocationByIP"], &locationsFound)
	if err != nil {
		return nil, err
	}
	if len(locationsFound) > 0 {
		return &locationsFound[0], nil
	}
	return nil, nil
}
