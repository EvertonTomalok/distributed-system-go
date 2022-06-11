package handlers

import (
	"net/http"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/methods"
	"github.com/gin-gonic/gin"
)

func GetAllMethods(c *gin.Context) {
	allMethods, err := methods.GetMethods(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var methodsResponse []dto.MethodResponse
	for _, method := range allMethods {
		methodsResponse = append(methodsResponse, dto.MethodResponse{Id: method.ID, Method: method.Name, Installments: int64(method.Installment)})
	}

	c.JSON(http.StatusOK, methodsResponse)
}
