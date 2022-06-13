package kafka

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/evertontomalok/distributed-system-go/internal/domain/broker"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
)

// The publisher will be injected when some command starts (server, wokers, orchestrator, etc.), by the Method NewPublisher
var Publisher message.Publisher

func PublishOrderMessageToTopic(topic string, order entities.Order, messageType string) error {
	msg, err := broker.NewOrderMessage(topic, order, messageType).Build()
	if err != nil {
		return err
	}

	return Publisher.Publish(topic, msg)
}

func PublishInternalMessageToTopic(topic string, internalMessage dto.BrokerInternalMessage) error {
	msg, err := broker.NewInternalMessage(internalMessage).Build()
	if err != nil {
		return err
	}
	return Publisher.Publish(topic, msg)
}
