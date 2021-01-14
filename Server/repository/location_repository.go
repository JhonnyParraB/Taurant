package repository

import (
	"encoding/json"

	"Taurant/driver"
	"Taurant/model"

	jsoniter "github.com/json-iterator/go"
)

//LocationRepository _
type LocationRepository interface {
	Create(location *model.Location) string
	FindByIP(ip string) *model.Location
	Update(uid string, location *model.Location) string
}

//LocationRepositoryDGraph _
type LocationRepositoryDGraph struct {
}

//Create _
func (b LocationRepositoryDGraph) Create(location *model.Location) error {
	location.UID = "_:" + location.IP
	err := driver.RunMutation(location)
	if err != nil {
		return err
	}
	locationFound, err := b.FindByIP(location.IP)
	if err != nil {
		return err
	}
	*(location) = *(locationFound)
	return nil
}

//Update _
func (b LocationRepositoryDGraph) Update(uid string, location *model.Location) error {
	location.UID = uid
	err := driver.RunMutation(location)
	if err != nil {
		return err
	}
	return nil
}

//FindByIP _
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
