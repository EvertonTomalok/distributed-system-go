package orders

import (
	"context"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/ports"
	uuid "github.com/satori/go.uuid"
)

var (
	OrdersDBAdapter ports.MethodsPort
)

func SaveOrder(ctx context.Context, orderRequest dto.OrderRequest) (string, error) {
	orderId := uuid.NewV4().String()
	return orderId, nil
}
