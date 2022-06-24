package handlers

import (
	"net/http"

	"github.com/evertontomalok/distributed-system-go/internal/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/core/methods"
	"github.com/evertontomalok/distributed-system-go/internal/infra/services/aws"
	"github.com/gin-gonic/gin"
)

// API Get Methods godoc
// @Summary Get Methods
// @Description Get available methods
// @Tags methods
// @Router /methods [get]
// @Accept json
// @Produce json
// @Success 200 {object} []dto.MethodResponse
// @Failure 500 "Something went wrong"
func GetAllMethods(c *gin.Context) {
	allMethods, err := methods.GetMethods(c.Request.Context())
	if err != nil {
		if err := c.AbortWithError(http.StatusNotFound, err); err != nil {
			aws.SendErrorToCloudWatch(c, err)
		}
		return
	}

	var methodsResponse []dto.MethodResponse
	for _, method := range allMethods {
		methodsResponse = append(methodsResponse, dto.MethodResponse{Id: method.ID, Method: method.Name, Installments: int64(method.Installment)})
	}

	c.JSON(http.StatusOK, methodsResponse)
}
