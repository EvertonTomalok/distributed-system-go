package dto

import "github.com/shopspring/decimal"

type OrderResponse struct {
	Value       decimal.Decimal `json:"value"`
	UserId      string          `json:"user_id"`
	Installment int8            `json:"installment"`
	Status      bool            `json:"status"`
	Method      string          `json:"method"`
}
