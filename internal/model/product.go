package model

import "github.com/shopspring/decimal"

type Product struct {
	ID          int             `json:"id"`
	Name        string          `json:"name`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
}
