package usecases

import (
	"net/http"

	"github.com/evertontomalok/distributed-system-go/internal/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/core/methods"
	"github.com/evertontomalok/distributed-system-go/internal/infra/services/aws"
	"github.com/gin-gonic/gin"
)

func GetAllMethods(c *gin.Context) *gin.Context {
	allMethods, err := methods.GetMethods(c.Request.Context())
	if err != nil {
		if err := c.AbortWithError(http.StatusNotFound, err); err != nil {
			aws.SendErrorToCloudWatch(c, err)
		}
		return c
	}

	var methodsResponse []dto.MethodResponse
	for _, method := range allMethods {
		methodsResponse = append(methodsResponse, dto.MethodResponse{Id: method.ID, Method: method.Name, Installments: int64(method.Installment)})
	}

	c.JSON(http.StatusOK, methodsResponse)
	return c
}
