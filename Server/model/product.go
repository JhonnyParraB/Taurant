package model

//Product type details
type Product struct {
	ID           string        `json:"product_id"`
	Name         string        `json:"product_name"`
	Price        int           `json:"price"`
	Transactions []Transaction `json:"has_transactions"`
}
