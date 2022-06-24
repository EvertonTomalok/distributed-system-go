package helpers

import (
	"github.com/evertontomalok/distributed-system-go/internal/core/domain/entities"
	"github.com/evertontomalok/distributed-system-go/internal/core/dto"
	kafkaAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
	"github.com/evertontomalok/distributed-system-go/pkg/broker"
)

func TriggerValidation(order entities.Order) error {
	if err := kafkaAdapter.PublishOrderMessageToTopic(broker.OrchestratorTopic, order, dto.StartEvent); err != nil {
		return err
	}
	return nil
}
