package entities

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID        string
	Value     decimal.Decimal
	MethodId  int64
	Method    Method
	UserId    string
	Status    sql.NullBool
	CreatedAt time.Time
	UpdatedAt time.Time
}
