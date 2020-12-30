package model

//Buyer type details
type Buyer struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Age          int           `json:"age"`
	Transacionts []Transaction `json:"perform"`
}
