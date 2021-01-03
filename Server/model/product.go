package model

//Product type details
type Product struct {
	UID           string          `predicate:"uid,omitempty"           endpoint:"uid,omitempty"`
	ID            string          `predicate:"product_id,omitempty"    endpoint:"id,omitempty"`
	Name          string          `predicate:"product_name,omitempty"  endpoint:"name,omitempty"`
	Price         int             `predicate:"price,omitempty"         endpoint:"price,omitempty"`
	ProductOrders *[]ProductOrder `predicate:"is_trade_in,omitempty"   endpoint:"product_orders,omitempty"`
}
