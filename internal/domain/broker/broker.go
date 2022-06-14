package broker

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	"github.com/oklog/ulid"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type OrderMessageBuilder struct {
	message *message.Message
	err     error
}

func UUID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now().UTC()), ulid.Monotonic(rand.New(rand.NewSource(time.Now().UTC().UnixNano())), 0)).String()
}

func NewOrderMessage(event string, order entities.Order, messageType string) *OrderMessageBuilder {
	json, err := json.Marshal(order)
	if err != nil {
		return &OrderMessageBuilder{err: err}
	}

	message := message.NewMessage(watermill.NewUUID(), json)
	message.Metadata.Set("aggregate_id", order.ID)
	message.Metadata.Set("event_id", UUID())
	message.Metadata.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()))
	message.Metadata.Set("message_type", messageType)

	return &OrderMessageBuilder{message: message}
}

func NewInternalMessage(brokerInternalMessage dto.BrokerInternalMessage) *OrderMessageBuilder {
	json, err := json.Marshal(brokerInternalMessage)
	if err != nil {
		return &OrderMessageBuilder{err: err}
	}

	message := message.NewMessage(watermill.NewUUID(), json)
	message.Metadata.Set("aggregate_id", brokerInternalMessage.ID)
	message.Metadata.Set("event_id", UUID())
	message.Metadata.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()))
	message.Metadata.Set("event", dto.ProcessInternalMessage)

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

func ParseOrderMessage(msg *message.Message) (dto.Order, error) {
	metadata := dto.Metadata{
		EventId:     msg.Metadata.Get("event_id"),
		AggregateId: msg.Metadata.Get("aggregate_id"),
		Timestamp:   msg.Metadata.Get("timestamp"),
		MessageType: msg.Metadata.Get("message_type"),
	}
	orderMessage := entities.Order{}
	err := json.Unmarshal(msg.Payload, &orderMessage)
	if err != nil {
		return dto.Order{}, err
	}

	return dto.Order{Metadata: metadata, Order: orderMessage}, nil
}

func ParseBrokerInternalMessage(msg *message.Message) (dto.BrokerInternalMessage, dto.Metadata, error) {
	metadata := dto.Metadata{
		EventId:     msg.Metadata.Get("event_id"),
		AggregateId: msg.Metadata.Get("aggregate_id"),
		Timestamp:   msg.Metadata.Get("timestamp"),
		Event:       msg.Metadata.Get("event"),
	}
	internalMessage := dto.BrokerInternalMessage{}
	err := json.Unmarshal(msg.Payload, &internalMessage)
	if err != nil {
		return dto.BrokerInternalMessage{}, metadata, err
	}

	return internalMessage, metadata, nil
}
