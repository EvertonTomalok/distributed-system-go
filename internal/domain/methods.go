package domain

import (
	"context"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/errors"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/ports"
)

var (
	MethodsDBAdapter ports.MethodsPort
)

func GetMethods(ctx context.Context) ([]entities.Method, error) {
	methods, err := MethodsDBAdapter.GetAllMethods(ctx)
	if err != nil {
		return nil, errors.InternalError
	}
	return methods, err
}
