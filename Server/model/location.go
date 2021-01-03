package model

//Location IP details

type Location struct {
	UID string `predicate:"uid,omitempty" endpoint:"uid,omitempty"`
	IP  string `predicate:"ip,omitempty"  endpoint:"ip,omitempty"`
}
