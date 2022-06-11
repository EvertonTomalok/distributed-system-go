package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	ordersRepository "github.com/evertontomalok/distributed-system-go/internal/domain/orders"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func PostOrder(c *gin.Context) {
	orderRequest := dto.OrderRequest{}

	err := c.ShouldBind(&orderRequest)

	if err == nil {
		orderId, err := ordersRepository.SaveOrder(c.Request.Context(), orderRequest)
		fmt.Printf(orderId, err)
		c.String(http.StatusOK, orderId)
		return
	}
	c.String(http.StatusNotModified, "hi")
}

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

	order := entities.Order{
		ID:       "e53c0252-1441-420c-8c01-9a22ae2005e2",
		UserId:   userId,
		Value:    value,
		Status:   true,
		MethodId: 1,
		Method: entities.Method{
			Name:        "credit_card",
			Installment: 1,
		},
	}

	orderResponse := dto.OrderResponse{
		Id:          order.ID,
		UserId:      order.UserId,
		Value:       order.Value,
		Installment: order.Method.Installment,
		Status:      order.Status,
		Method:      order.Method.Name,
	}
	c.JSON(http.StatusOK, orderResponse)
}
