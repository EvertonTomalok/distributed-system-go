package dto

import (
	"github.com/shopspring/decimal"
)

type OrderResponse struct {
	Id          string          `json:"id"`
	Value       decimal.Decimal `json:"value"`
	UserId      string          `json:"user_id"`
	Installment int64           `json:"installment"`
	Status      string          `json:"status" swaggertype:"string"`
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

type OrdersResponse struct {
	Offset int64           `json:"offset"`
	Limit  int64           `json:"limit"`
	Total  int             `json:"total"`
	Orders []OrderResponse `json:"orders"`
}
