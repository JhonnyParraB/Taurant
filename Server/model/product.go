package model

//Product type details
type Product struct {
	UID          string        `predicate:"uid,omitempty"`
	ID           string        `predicate:"product_id,omitempty"`
	Name         string        `predicate:"product_name,omitempty"`
	Price        int           `predicate:"price,omitempty"`
	Transactions []Transaction `predicate:"has_transactions,omitempty"`
}
