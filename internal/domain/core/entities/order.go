package entities

import "github.com/shopspring/decimal"

type Order struct {
	ID         string
	Value      decimal.Decimal
	MethodId   int64
	Method     Method
	UserId     string
	Status     bool
	created_at string
	expire_at  string
}
