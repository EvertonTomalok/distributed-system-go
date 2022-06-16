package orchestrator

import (
	"context"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/evertontomalok/distributed-system-go/internal/app"
	"github.com/evertontomalok/distributed-system-go/internal/app/utils"
	"github.com/evertontomalok/distributed-system-go/internal/domain/broker"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	eventSource "github.com/evertontomalok/distributed-system-go/internal/domain/events"
	log "github.com/sirupsen/logrus"

	kafkaAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
)

func StartOrchestrator(ctx context.Context, config app.Config) {
	subscriber := kafkaAdapter.NewSubscriber("orchestrator", config.Kafka.Host, config.Kafka.Port)
	kafkaAdapter.Publisher = kafkaAdapter.NewPublisher(config.Kafka.Host, config.Kafka.Port)

	done := utils.MakeDoneSignal()

	messages, err := subscriber.Subscribe(context.Background(), broker.OrchestatratorTopic)
	if err != nil {
		log.Panicf("Trying to start orchestrator, some error ocurred: %+v: ", err)
	}
	go process(messages)

	<-done
	subscriber.Close()
}

func process(messages <-chan *message.Message) {
	for msg := range messages {
		if err := processMessage(msg); err != nil {
			log.Errorf("Something went wrong trying to process message %+v | err: %+v", msg, err)
			// TODO send this message to a dead letter
		}
		msg.Ack()
	}
}

func processMessage(msg *message.Message) error {
	event := msg.Metadata.Get("event")

	switch event {
	case dto.StartEvent:
		triggerWorkers(msg)
	case dto.ResultValidateBalance:
		updateStep(msg, "Balance Validated")
	case dto.ResultValidateUserStatus:
		updateStep(msg, "Status Validated")
	default:
		updateStep(msg, "Default")
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

	err := eventSource.CreateEventSource(context.Background(), internalMessage)

	if err != nil {
		log.Printf("Some error ocurred trying to save event source: %+v", err)
		return
	}

	for _, topic := range [2]string{broker.UserStatusValidatorTopic, broker.UserBalanceValidatorTopic} {
		wg.Add(1)
		go func(t string, i dto.BrokerInternalMessage) {
			defer wg.Done()
			kafkaAdapter.PublishInternalMessageToTopic(t, i, dto.StartEvent)
		}(topic, internalMessage)
	}
	wg.Wait()
}

func updateStep(msg *message.Message, message string) {
	internalMessage, metadata, err := broker.ParseBrokerInternalMessage(msg)
	step := dto.EventSteps{
		Event:   metadata.Event,
		Status:  true,
		Message: message,
	}

	err = eventSource.UpdateStep(context.Background(), internalMessage.ID, step)

	if err != nil {
		log.Printf("Some error ocurred trying to update event source: %+v", err)
		return
	}

}
