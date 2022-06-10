package domain

import "github.com/shopspring/decimal"

type Order struct {
	Value  decimal.Decimal
	Method Method
	UserId string
	Status bool
}
