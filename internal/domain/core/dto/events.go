package dto

import (
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	"github.com/shopspring/decimal"
)

const (
	StartEvent               string = "start_event"
	ProcessInternalMessage          = "process_internal_message"
	ResultValidateUserStatus        = "result_validate_user_status"
	ResultValidateBalance           = "result_validate_balance"
)

type EventSteps struct {
	Event   string `bson:"event"`
	Message string `bson:"message"`
	Status  bool   `bson:"status"`
}

type BrokerInternalMessage struct {
	ID           string          `json:"id"`
	Value        decimal.Decimal `json:"value"`
	MethodId     int64           `json:"method_id"`
	Method       string          `json:"method"`
	Installments int64           `json:"installments"`
	UserId       string          `json:"user_id"`
	Status       bool            `json:"status"`
	Steps        []EventSteps    `json:"steps"`
}

type OrderMessage struct {
	Status       string          `json:"status"`
	OrderId      string          `json:"order_id"`
	UserId       string          `json:"user_id"`
	Method       string          `json:"type"`
	Instalmments int64           `json:"installments"`
	Value        decimal.Decimal `json:"value"`
}

type Metadata struct {
	EventId     string
	AggregateId string
	Timestamp   string
	MessageType string
	Event       string
}

type Order struct {
	Metadata
	Order entities.Order
}
