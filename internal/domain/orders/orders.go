package orders

import (
	"context"
	"time"

	"github.com/evertontomalok/distributed-system-go/internal/domain/broker"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/errors"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/ports"
	"github.com/evertontomalok/distributed-system-go/internal/domain/methods"
	kafkaAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
	uuid "github.com/satori/go.uuid"
)

// The adapter will be injected when the command starts the program. See the injection in
// github.com/evertontomalok/distributed-system-go/internal/app/config::InitDb()
var (
	OrdersDBAdapter ports.OrdersPort
)

func SaveOrder(ctx context.Context, orderRequest dto.OrderRequest) (string, error) {
	method, err := methods.MethodsDBAdapter.GetMethodByNameAndInstallment(ctx, orderRequest.Method, orderRequest.Installment)
	if err != nil {
		return "", errors.InvalidMethod
	}

	now := time.Now()

	orderUUID := uuid.NewV4().String()
	order := entities.Order{
		ID:        orderUUID,
		Value:     orderRequest.Value,
		MethodId:  method.ID,
		UserId:    orderRequest.UserId,
		CreatedAt: now,
		UpdatedAt: now,
	}

	orderId, err := OrdersDBAdapter.PostOrder(ctx, order)

	if err != nil {
		return "", errors.InvalidOrder
	}

	for _, topic := range [2]string{broker.UserStatusValidatorTopic, broker.UserBalanceValidatorTopic} {
		kafkaAdapter.PublishOrderMessageToTopic(topic, order)
	}

	return orderId, nil
}
