package dto

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type OrderResponse struct {
	Id          string          `json:"id"`
	Value       decimal.Decimal `json:"value"`
	UserId      string          `json:"user_id"`
	Installment int64           `json:"installment"`
	Status      sql.NullBool    `json:"status" swaggertype:"boolean"`
	Method      string          `json:"method"`
}

type OrderRequest struct {
	Value       decimal.Decimal `form:"value" binding:"required"`
	UserId      string          `form:"user_id" binding:"required"`
	Installment int64           `form:"installment" binding:"required"`
	Method      string          `form:"method" binding:"required"`
}

type MethodResponse struct {
	Id           int64  `json:"id"`
	Method       string `json:"method"`
	Installments int64  `json:"installments"`
}
