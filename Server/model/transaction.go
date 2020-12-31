package model

//Transaction type details
type Transaction struct {
	UID      string    `predicate:"uid,omitempty"`
	ID       string    `predicate:"transaction_id,omitempty"`
	Buyer    Buyer     `predicate:"is_made_by,omitempty"`
	IP       string    `predicate:"ip,omitempty"`
	Device   string    `predicate:"device,omitempty"`
	Products []Product `predicate:"trade,omitempty"`
}
