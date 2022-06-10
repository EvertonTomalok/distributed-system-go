package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/evertontomalok/distributed-system-go/internal/domain"
	"github.com/evertontomalok/distributed-system-go/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func GetOrdersByUserId(c *gin.Context) {
	userId := c.Param("userId")
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 100
	}

	fmt.Println(offset, limit)

	value, _ := decimal.NewFromString("10.00")

	order := domain.Order{
		UserId: userId,
		Value:  value,
		Status: true,
		Method: domain.Method{
			Name:        "credit_card",
			Installment: 1,
		},
	}

	orderResponse := dto.OrderResponse{
		UserId:      order.UserId,
		Value:       order.Value,
		Installment: order.Method.Installment,
		Status:      order.Status,
		Method:      order.Method.Name,
	}
	c.JSON(http.StatusOK, orderResponse)
}
