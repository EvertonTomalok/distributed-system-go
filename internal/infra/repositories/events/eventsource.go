package event

import (
	"context"

	"github.com/evertontomalok/distributed-system-go/internal/app/ports"
	"github.com/evertontomalok/distributed-system-go/internal/core/dto"
)

var EventsAdapter ports.EventsSourcePort

func CreateEventSource(ctx context.Context, brokerInternalMessage dto.BrokerInternalMessage) error {
	err := EventsAdapter.SaveEvent(ctx, brokerInternalMessage)

	if err != nil {
		return err
	}
	return nil
}

func UpdateStep(ctx context.Context, orderId string, step dto.EventSteps) error {
	err := EventsAdapter.UpdateEventStep(ctx, orderId, step)

	if err != nil {
		return err
	}
	return nil
}
