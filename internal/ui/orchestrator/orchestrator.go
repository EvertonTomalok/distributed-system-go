package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/evertontomalok/distributed-system-go/internal/app"
	"github.com/evertontomalok/distributed-system-go/internal/core/domain/entities"
	"github.com/evertontomalok/distributed-system-go/internal/core/dto"
	eventSourceRepository "github.com/evertontomalok/distributed-system-go/internal/infra/repositories/events"
	"github.com/evertontomalok/distributed-system-go/internal/infra/services/aws"
	"github.com/evertontomalok/distributed-system-go/pkg/broker"
	"github.com/evertontomalok/distributed-system-go/pkg/utils"
	log "github.com/sirupsen/logrus"

	kafkaAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
	ordersRepository "github.com/evertontomalok/distributed-system-go/internal/infra/repositories/orders"
)

func StartOrchestrator(ctx context.Context, config app.Config) {
	subscriber := kafkaAdapter.NewSubscriber("orchestrator", config.Kafka.Host, config.Kafka.Port)
	kafkaAdapter.Publisher = kafkaAdapter.NewPublisher(config.Kafka.Host, config.Kafka.Port)

	done := utils.MakeDoneSignal()

	messages, err := subscriber.Subscribe(context.Background(), broker.OrchestratorTopic)
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
			log.Errorf("Something went wrong trying to process message %+v | err: %+v\n\n\n\n", msg, err)
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
	case dto.CompensationValidateUserStatus:
		updateStep(msg, "User is invalid.")
	case dto.CompensationBalanceStatus:
		updateStep(msg, "Balance is invalid.")
	default:
		msgReceived := fmt.Sprintf("Message event received, and it can't be processed -> %+v", msg)
		err := errors.New(msgReceived)
		aws.SendErrorToCloudWatch(context.TODO(), err)
	}

	if err := orderIsCompleted(context.Background(), msg); err != nil {
		return err
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

	err := eventSourceRepository.CreateEventSource(msg.Context(), internalMessage)

	if err != nil {
		aws.SendErrorToCloudWatch(msg.Context(), err)
		return
	}

	for _, topic := range [2]string{broker.UserStatusValidatorTopic, broker.UserBalanceValidatorTopic} {
		wg.Add(1)
		go func(t string, i dto.BrokerInternalMessage) {
			defer wg.Done()
			err = kafkaAdapter.PublishInternalMessageToTopic(t, i, dto.StartEvent)
			if err != nil {
				aws.SendErrorToCloudWatch(msg.Context(), err)
			}
		}(topic, internalMessage)
	}
	wg.Wait()
}

func updateStep(msg *message.Message, message string) {
	internalMessage, metadata, err := broker.ParseBrokerInternalMessage(msg)
	if err != nil {
		aws.SendErrorToCloudWatch(msg.Context(), err)
		return
	}
	step := dto.EventSteps{
		Event:   metadata.Event,
		Status:  internalMessage.Status,
		Message: message,
	}

	err = eventSourceRepository.UpdateStep(context.Background(), internalMessage.ID, step)

	if err != nil {
		aws.SendErrorToCloudWatch(msg.Context(), err)
		return
	}

	if !internalMessage.Status && (metadata.Event == dto.ResultValidateBalance || metadata.Event == dto.ResultValidateUserStatus) {
		step := dto.EventSteps{
			Event:   dto.CompensationStarted,
			Status:  true,
			Message: fmt.Sprintf("Compensation started to %s", metadata.Event),
		}
		err = eventSourceRepository.UpdateStep(context.Background(), internalMessage.ID, step)
		if err != nil {
			aws.SendErrorToCloudWatch(msg.Context(), err)
			return
		}
		compensationTrigger(internalMessage, &metadata)
	}
}

func compensationTrigger(internalMessage dto.BrokerInternalMessage, metadata *dto.Metadata) {
	switch metadata.Event {
	case dto.ResultValidateBalance:
		err := kafkaAdapter.PublishInternalMessageToTopic(broker.UserBalanceCompensationTopic, internalMessage, dto.CompensationBalanceStatus)
		if err != nil {
			aws.SendErrorToCloudWatch(context.TODO(), err)
		}
	case dto.ResultValidateUserStatus:
		err := kafkaAdapter.PublishInternalMessageToTopic(broker.UserStatusCompensationTopic, internalMessage, dto.CompensationValidateUserStatus)
		if err != nil {
			aws.SendErrorToCloudWatch(context.TODO(), err)
			return
		}
	}
}

func orderIsCompleted(ctx context.Context, msg *message.Message) error {
	internalMessage, _, _ := broker.ParseBrokerInternalMessage(msg)
	if internalMessage.ID == "" {
		return nil
	}
	doc, e := eventSourceRepository.EventsAdapter.GetDocumentByOrderId(msg.Context(), internalMessage.ID)

	if e != nil {
		return e
	}

	if len(doc.Steps) < 3 {
		return nil
	}

	if allStepsOk(doc.Steps) {
		err := ordersRepository.OrdersDBAdapter.UpdateStatusByOrderId(ctx, internalMessage.ID, entities.APPROVED)
		return err
	}
	err := ordersRepository.OrdersDBAdapter.UpdateStatusByOrderId(ctx, internalMessage.ID, entities.CANCELED)
	return err
}

func allStepsOk(steps []dto.EventSteps) bool {
	// Todo implement a better logic to these steps validation
	for _, step := range steps {
		if !step.Status {
			return false
		}
	}
	return true
}
