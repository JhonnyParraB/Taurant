package model

//Transaction type details
type Transaction struct {
	ID       string    `json:"id"`
	Buyer    Buyer     `json:"is_made_by"`
	IP       string    `json:"ip"`
	Device   string    `json:"device"`
	Products []Product `json:"trade"`
}
