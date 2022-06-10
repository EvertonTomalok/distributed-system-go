package entities

import "github.com/shopspring/decimal"

type Order struct {
	ID       int64
	Value    decimal.Decimal
	MethodId int64
	Method   Method
	UserId   string
	Status   bool
}
