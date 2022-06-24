package handlers

import (
	"github.com/evertontomalok/distributed-system-go/internal/app/usecases"
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
	usecases.GetAllMethods(c)
}
