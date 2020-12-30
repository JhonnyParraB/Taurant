package model

//Transaction type details
type Transaction struct {
	ID       string    `json:"id"`
	Buyer    Buyer     `json:"buyer"`
	IP       string    `json:"ip"`
	Device   string    `json:"device"`
	Products []Product `json:"products"`
}
