package ports

import (
	"context"
	"database/sql"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
)

type MethodsPort interface {
	GetAllMethods(ctx context.Context) ([]entities.Method, error)
}

type DB interface {
	DB() *sql.DB
}
