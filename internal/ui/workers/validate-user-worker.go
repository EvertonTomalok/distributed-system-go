package workers

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/evertontomalok/distributed-system-go/internal/app"
	"github.com/evertontomalok/distributed-system-go/internal/app/utils"
	"github.com/evertontomalok/distributed-system-go/internal/domain/broker"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"

	kafkaAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
	userapi "github.com/evertontomalok/distributed-system-go/internal/infra/services/user-api"
)

func StartValidateUserStatus(ctx context.Context, config app.Config) {
	router := kafkaAdapter.NewRouter()
	subscriber := kafkaAdapter.NewSubscriber("user-status-consumer", config.Kafka.Host, config.Kafka.Port)
	kafkaAdapter.Publisher = kafkaAdapter.NewPublisher(config.Kafka.Host, config.Kafka.Port)

	router.AddNoPublisherHandler(
		"validate-user-status-orders",
		broker.UserStatusValidatorTopic,
		subscriber,
		validateUserStatusOrder,
	)
	router.AddNoPublisherHandler(
		"compensate-user-status-orders",
		broker.UserStatusCompensationTopic,
		subscriber,
		rowBackUserStatusOrder,
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

func validateUserStatusOrder(msg *message.Message) error {
	internalMessage, metadata, err := broker.ParseBrokerInternalMessage(msg)
	if err != nil {
		log.Printf("validate user status -> %+v | %+v | %+v \n\n", internalMessage, metadata, err)
	}

	userStatusResponse := userapi.UserAdapter.GetUserStatus(internalMessage.UserId)
	if userStatusResponse.IsValid == true {
		kafkaAdapter.PublishInternalMessageToTopic(broker.OrchestratorTopic, internalMessage, dto.ResultValidateUserStatus)
	} else {
		internalMessage.Status = false
		kafkaAdapter.PublishInternalMessageToTopic(broker.UserStatusCompensationTopic, internalMessage, dto.ResultValidateUserStatus)
	}

	return nil
}

func rowBackUserStatusOrder(msg *message.Message) error {
	internalMessage, metadata, err := broker.ParseBrokerInternalMessage(msg)
	if err != nil {
		log.Printf("%+v | %+v | %+v \n\n", internalMessage, metadata, err)
	}

	kafkaAdapter.PublishInternalMessageToTopic(broker.OrchestratorTopic, internalMessage, dto.CompensationValidateUserStatus)
	return nil
}
