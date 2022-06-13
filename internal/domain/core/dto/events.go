package dto

import (
	"github.com/shopspring/decimal"
)

const (
	StartEvent               string = "start_event"
	ValidateUserStatus              = "validate_user_status"
	ResultValidateUserStatus        = "result_validate_user_status"
	ValidateBalance                 = "validate_balance"
	ResultValidateBalance           = "result_validate_balance"
)

type BrokerInternalMessage struct {
	ID           string          `json:"id"`
	Value        decimal.Decimal `json:"value"`
	MethodId     int64           `json:"method_id"`
	Method       string          `json:"method"`
	Installments int64           `json:"installments"`
	UserId       string          `json:"user_id"`
	Status       bool            `json:"status"`
	Metadata     string          `json:"metadata"`
}
