package ports

import (
	"context"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
)

type MethodsPort interface {
	GetAllMethods(ctx context.Context) ([]entities.Method, error)
}
