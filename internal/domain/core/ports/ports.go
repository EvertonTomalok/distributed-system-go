package ports

import (
	"context"
	"database/sql"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
)

type MethodsPort interface {
	GetAllMethods(ctx context.Context) ([]entities.Method, error)
	GetMethodByNameAndInstallment(ctx context.Context, methodName string, installments int64) (entities.Method, error)
}

type OrdersPort interface {
	PostOrder(ctx context.Context, orderRequest entities.Order) (string, error)
	GetOrdersByUserId(ctx context.Context, userId string, offset int64, limit int64) ([]entities.Order, error)
}

type DB interface {
	DB() *sql.DB
}
