package model

//ProductOrder type details
type ProductOrder struct {
	UID         string       `predicate:"uid,omitempty"        endpoint:"uid,omitempty"`
	Product     *Product     `predicate:"trade,omitempty"      endpoint:"product,omitempty"`
	Quantity    int          `predicate:"quantity,omitempty"   endpoint:"quantity,omitempty"`
	Transaction *Transaction `predicate:"belong_to,omitempty"  endpoint:"transaction,omitempty"`
}
