package model

//Transaction type details
type Transaction struct {
	UID           string          `predicate:"uid,omitempty"              endpoint:"uid,omitempty"`
	ID            string          `predicate:"transaction_id,omitempty"   endpoint:"id,omitempty"`
	Buyer         *Buyer          `predicate:"is_made_by,omitempty"       endpoint:"buyer,omitempty"`
	Date          int64           `predicate:"transaction_date,omitempty" endpoint:"date,omitempty"`
	Location      *Location       `predicate:"location,omitempty"         endpoint:"location,omitempty"`
	Device        string          `predicate:"device,omitempty"           endpoint:"device,omitempty"`
	ProductOrders *[]ProductOrder `predicate:"include,omitempty"          endpoint:"product_orders,omitempty"`
}
