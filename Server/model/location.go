package model

//Location IP details

type Location struct {
	UID string `predicate:"uid,omitempty"`
	IP  string `predicate:"ip,omitempty"`
}
