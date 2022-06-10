package ports

import (
	"context"
	"database/sql"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
)

type MethodsPort interface {
	GetAllMethods(ctx context.Context) ([]entities.Method, error)
}

type OrdersPort interface {
	PostOrder(ctx context.Context, orderRequest dto.OrderRequest) (string, error)
}

type DB interface {
	DB() *sql.DB
}
