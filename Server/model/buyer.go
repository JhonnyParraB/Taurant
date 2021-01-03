package model

//Buyer type details
type Buyer struct {
	UID          string        `predicate:"uid,omitempty"`
	ID           string        `predicate:"buyer_id,omitempty"    endpoint:"id,omitempty"`
	Name         string        `predicate:"buyer_name,omitempty"  endpoint:"name,omitempty"`
	Age          int           `predicate:"age,omitempty"         endpoint:"age,omitempty"`
	Transactions []Transaction `predicate:"perform,omitempty"     endpoint:"transactions,omitempty"`
	SharedIP     string        `predicate:"shared_ip,omitempty"   endpoint:"sharedIP,omitempty"`
}
