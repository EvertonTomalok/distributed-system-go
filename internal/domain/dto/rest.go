package dto

import "github.com/shopspring/decimal"

type OrderResponse struct {
	Id          string          `json:"id"`
	Value       decimal.Decimal `json:"value"`
	UserId      string          `json:"user_id"`
	Installment int8            `json:"installment"`
	Status      bool            `json:"status"`
	Method      string          `json:"method"`
}

type OrderRequest struct {
	Value       decimal.Decimal `form:"value" binding:"required"`
	UserId      string          `form:"user_id" binding:"required"`
	Installment int8            `form:"installment" binding:"required"`
	Method      string          `form:"method" binding:"required"`
}
