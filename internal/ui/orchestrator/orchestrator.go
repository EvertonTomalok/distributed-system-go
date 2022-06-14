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
	event "github.com/evertontomalok/distributed-system-go/internal/domain/events"

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
	messageType := msg.Metadata.Get("message_type")
	switch messageType {
	case dto.StartEvent:
		triggerWorkers(msg)
	case dto.ResultValidateBalance:
		updateStep(msg)
	case dto.ResultValidateUserStatus:
		updateStep(msg)
	default:
		updateStep(msg)
	}

	return nil
}

func triggerWorkers(msg *message.Message) {
	message, _ := broker.ParseOrderMessage(msg)
	var wg sync.WaitGroup

	step := dto.EventSteps{
		Event:   dto.StartEvent,
		Status:  true,
		Message: "Started with success",
	}

	steps := make([]dto.EventSteps, 0)
	steps = append(steps, step)

	internalMessage := dto.BrokerInternalMessage{
		ID:           message.Order.ID,
		Value:        message.Order.Value,
		MethodId:     message.Order.MethodId,
		Method:       message.Order.Method.Name,
		Installments: int64(message.Order.Method.Installment),
		UserId:       message.Order.UserId,
		Status:       true,
		Steps:        steps,
	}

	err := event.CreateEventSource(context.Background(), internalMessage)

	if err != nil {
		log.Printf("Some error ocurred trying to save event source: %+v", err)
		return
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

func updateStep(msg *message.Message) {
	internalMessage, metadata, err := broker.ParseBrokerInternalMessage(msg)
	step := dto.EventSteps{
		Event:   metadata.Event,
		Status:  true,
		Message: metadata.MessageType,
	}

	err = event.UpdateStep(context.Background(), internalMessage.ID, step)

	if err != nil {
		log.Printf("Some error ocurred trying to update event source: %+v", err)
		return
	}

}
