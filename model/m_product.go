package model

type Products []Product
type Product struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Price  string `json:"price"`
	Active int    `json:"active"`
}

func (Product) TableName() string {
	return "product"
}
