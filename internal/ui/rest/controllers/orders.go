package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/evertontomalok/distributed-system-go/internal/domain/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func PostOrder(c *gin.Context) {
	order := dto.OrderRequest{}

	err := c.ShouldBind(&order)

	if err == nil {
		fmt.Println(order.Installment)
		fmt.Println(order.Method)
		fmt.Println(order.Value)
		fmt.Println(order.UserId)
	} else {
		fmt.Println(err)
	}

	c.String(http.StatusOK, "Ok")
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
		UserId: userId,
		Value:  value,
		Status: true,
		Method: entities.Method{
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
