package methods

import (
	"context"

	"github.com/evertontomalok/distributed-system-go/internal/core/domain/entities"
	"github.com/evertontomalok/distributed-system-go/internal/core/domain/errors"
	"github.com/evertontomalok/distributed-system-go/internal/core/domain/ports"
)

// The adapter will be injected when the command starts the program. See the injection in
// github.com/evertontomalok/distributed-system-go/internal/app/config::InitDb()
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
