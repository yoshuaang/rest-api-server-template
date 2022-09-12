package model

type Products []Product
type Product struct {
	ID     int
	Name   string
	Price  string
	Active int
}

func (Product) TableName() string {
	return "product"
}
