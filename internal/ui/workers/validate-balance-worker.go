package workers

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/evertontomalok/distributed-system-go/internal/app"
	"github.com/evertontomalok/distributed-system-go/internal/core/dto"
	"github.com/evertontomalok/distributed-system-go/pkg/broker"
	"github.com/evertontomalok/distributed-system-go/pkg/utils"

	kafkaAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
)

func StartValidateBalance(ctx context.Context, config app.Config) {
	router := kafkaAdapter.NewRouter()
	subscriber := kafkaAdapter.NewSubscriber("balance-consumer", config.Kafka.Host, config.Kafka.Port)
	kafkaAdapter.Publisher = kafkaAdapter.NewPublisher(config.Kafka.Host, config.Kafka.Port)

	router.AddNoPublisherHandler(
		"validate-balance-orders",
		broker.UserBalanceValidatorTopic,
		subscriber,
		validateBalanceOrder,
	)
	router.AddNoPublisherHandler(
		"compensate-balance-orders",
		broker.UserBalanceCompensationTopic,
		subscriber,
		rowBackBalanceOrder,
	)
	done := utils.MakeDoneSignal()
	go func() {
		if err := router.Run(ctx); err != nil {
			log.Panicf("%+v\n\n", err)
		}
	}()
	<-done
	router.Close()
}

func validateBalanceOrder(msg *message.Message) error {
	internalMessage, metadata, err := broker.ParseBrokerInternalMessage(msg)
	if err != nil {
		log.Printf("%+v | %+v | %+v \n\n", internalMessage, metadata, err)
	}
	v, _ := internalMessage.Value.Float64()
	if v <= 10000.00 {
		if err := kafkaAdapter.PublishInternalMessageToTopic(broker.OrchestratorTopic, internalMessage, dto.ResultValidateBalance); err != nil {
			return err
		}
	} else {
		internalMessage.Status = false
		if err := kafkaAdapter.PublishInternalMessageToTopic(broker.OrchestratorTopic, internalMessage, dto.CompensationBalanceStatus); err != nil {
			return err
		}
	}
	return nil
}

func rowBackBalanceOrder(msg *message.Message) error {
	internalMessage, metadata, err := broker.ParseBrokerInternalMessage(msg)
	if err != nil {
		log.Printf("validate balance -> %+v | %+v | %+v \n\n", internalMessage, metadata, err)
	}

	if err := kafkaAdapter.PublishInternalMessageToTopic(broker.OrchestratorTopic, internalMessage, dto.CompensationBalanceStatus); err != nil {
		return err
	}
	return nil
}
