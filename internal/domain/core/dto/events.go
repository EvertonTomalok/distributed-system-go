package dto

import (
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	StartEvent                     string = "start_event"
	ProcessInternalMessage                = "process_internal_message"
	CompensationStarted                   = "compensation_started"
	CompensationBalanceStatus             = "compensation_balance_status"
	ResultValidateUserStatus              = "result_validate_user_status"
	CompensationValidateUserStatus        = "compensation_validate_user_status"
	ResultValidateBalance                 = "result_validate_balance"
)

type EventSteps struct {
	Event   string `bson:"event"`
	Message string `bson:"message"`
	Status  bool   `bson:"status"`
}

type BrokerInternalMessage struct {
	ID           string          `json:"id" bson:"id"`
	Value        decimal.Decimal `json:"value" bson:"value"`
	MethodId     int64           `json:"method_id" bson:"method_id"`
	Method       string          `json:"method" bson:"method"`
	Installments int64           `json:"installments" bson:"installments"`
	UserId       string          `json:"user_id" bson:"user_id"`
	Status       bool            `json:"status" bson:"status"`
	Steps        []EventSteps    `json:"steps" bson:"steps"`
}

type Document struct {
	ID           string               `json:"id" bson:"id"`
	Value        primitive.Decimal128 `json:"value" bson:"value"`
	MethodId     int64                `json:"method_id" bson:"method_id"`
	Method       string               `json:"method" bson:"method"`
	Installments int64                `json:"installments" bson:"installments"`
	UserId       string               `json:"user_id" bson:"user_id"`
	Status       bool                 `json:"status" bson:"status"`
	Steps        []EventSteps         `json:"steps" bson:"steps"`
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
