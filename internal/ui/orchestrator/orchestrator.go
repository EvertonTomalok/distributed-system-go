package orchestrator

import (
	"context"
	"log"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/evertontomalok/distributed-system-go/internal/app"
	"github.com/evertontomalok/distributed-system-go/internal/app/utils"
	"github.com/evertontomalok/distributed-system-go/internal/domain/broker"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"

	kafkaAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
)

func StartOrchestrator(ctx context.Context, config app.Config) {
	router := kafkaAdapter.NewRouter()
	subscriber := kafkaAdapter.NewSubscriber("orchestrator", config.Kafka.Host, config.Kafka.Port)
	kafkaAdapter.Publisher = kafkaAdapter.NewPublisher(config.Kafka.Host, config.Kafka.Port)

	router.AddNoPublisherHandler(
		"orchestrator",
		broker.OrchestatratorTopic,
		subscriber,
		processMessage,
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

func processMessage(msg *message.Message) error {
	payload := string(msg.Payload)
	log.Printf("received message: %s, payload: %s, metadata: %+v", msg.UUID, payload, msg.Metadata)
	messageType := msg.Metadata.Get("message_type")
	switch messageType {
	case dto.StartEvent:
		triggerWorkers(msg)
	}

	return nil
}

func triggerWorkers(msg *message.Message) {
	message, _ := broker.ParseOrderMessage(msg)
	var wg sync.WaitGroup

	internalMessage := dto.BrokerInternalMessage{
		ID:           message.Order.ID,
		Value:        message.Order.Value,
		MethodId:     message.Order.MethodId,
		Installments: int64(message.Order.Method.Installment),
		UserId:       message.Order.UserId,
	}

	for _, topic := range [2]string{broker.UserStatusValidatorTopic, broker.UserBalanceValidatorTopic} {
		wg.Add(1)

		go func(t string, i dto.BrokerInternalMessage) {
			defer wg.Done()
			kafkaAdapter.PublishInternalMessageToTopic(t, i)
		}(topic, internalMessage)
	}
	wg.Wait()
}
