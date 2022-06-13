package workers

import (
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/evertontomalok/distributed-system-go/internal/app"
	"github.com/evertontomalok/distributed-system-go/internal/app/utils"
	"github.com/evertontomalok/distributed-system-go/internal/domain/broker"

	kafkaAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
)

func StartValidateBalance(ctx context.Context, config app.Config) {
	router := kafkaAdapter.NewRouter()
	subscriber := kafkaAdapter.NewSubscriber("balance-consumer", config.Kafka.Host, config.Kafka.Port)

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
		log.Println("Worker Started!")
		if err := router.Run(ctx); err != nil {
			log.Panicf("%+v\n\n", err)
		}
	}()
	<-done
	router.Close()
}

func validateBalanceOrder(msg *message.Message) error {
	log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
	internalMessage, metadata, err := broker.ParseBrokerInternalMessage(msg)
	log.Printf("%+v | %+v | %+v \n\n", internalMessage, metadata, err)
	return nil
}

func rowBackBalanceOrder(msg *message.Message) error {
	log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
	internalMessage, metadata, err := broker.ParseBrokerInternalMessage(msg)
	log.Printf("%+v | %+v | %+v \n\n", internalMessage, metadata, err)
	return nil
}
