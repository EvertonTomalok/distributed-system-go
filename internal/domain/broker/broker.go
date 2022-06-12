package broker

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	"github.com/oklog/ulid"
	"github.com/shopspring/decimal"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

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
}

type Order struct {
	Metadata
	Message OrderMessage
}

type OrderMessageBuilder struct {
	message *message.Message
	err     error
}

func UUID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now().UTC()), ulid.Monotonic(rand.New(rand.NewSource(time.Now().UTC().UnixNano())), 0)).String()
}

func NewOrderMessage(event string, order entities.Order) *OrderMessageBuilder {
	json, err := json.Marshal(order)
	if err != nil {
		return &OrderMessageBuilder{err: err}
	}

	message := message.NewMessage(watermill.NewUUID(), json)
	message.Metadata.Set("aggregate_id", order.ID)
	message.Metadata.Set("event_id", UUID())
	message.Metadata.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()))

	return &OrderMessageBuilder{message: message}
}

func (m *OrderMessageBuilder) WithAggregate(aggregateType string) *OrderMessageBuilder {
	m.message.Metadata.Set("aggregate_id", UUID())
	m.message.Metadata.Set("aggregate_type", aggregateType)

	return m
}

func (m *OrderMessageBuilder) Build() (*message.Message, error) {
	return m.message, m.err
}

func ParseOrderMessage(msg *message.Message) (Order, error) {
	metadata := Metadata{
		EventId:     msg.Metadata.Get("event_id"),
		AggregateId: msg.Metadata.Get("aggregate_id"),
		Timestamp:   msg.Metadata.Get("timestamp"),
	}
	orderMessage := OrderMessage{}
	err := json.Unmarshal(msg.Payload, &orderMessage)
	if err != nil {
		return Order{}, err
	}

	return Order{Metadata: metadata, Message: orderMessage}, nil
}
