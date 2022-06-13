package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	PROCESSING string = "PROCESSING"
	APPROVED          = "APPROVED"
	CANCELED          = "CANCELED"
)

type Order struct {
	ID        string
	Value     decimal.Decimal
	MethodId  int64
	Method    Method
	UserId    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
