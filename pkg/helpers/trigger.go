package helpers

import (
	"github.com/evertontomalok/distributed-system-go/internal/domain/broker"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	kafkaAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
)

func TriggerValidation(order entities.Order) error {
	if err := kafkaAdapter.PublishOrderMessageToTopic(broker.OrchestratorTopic, order, dto.StartEvent); err != nil {
		return err
	}
	return nil
}
