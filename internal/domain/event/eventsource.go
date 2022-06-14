package event

import (
	"context"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/ports"
)

var EventsAdapter ports.EventsSourcePort

func CreateEventSource(ctx context.Context, brokerInternalMessage dto.BrokerInternalMessage) error {
	err := EventsAdapter.SaveEvent(ctx, brokerInternalMessage)

	if err != nil {
		return err
	}
	return nil
}
