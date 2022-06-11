package handlers

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/evertontomalok/distributed-system-go/internal/domain/broker"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/errors"
	ordersRepository "github.com/evertontomalok/distributed-system-go/internal/domain/orders"
	kafkaAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
	"github.com/gin-gonic/gin"
)

func PostOrder(c *gin.Context) {
	orderRequest := dto.OrderRequest{}

	err := c.ShouldBind(&orderRequest)

	if err == nil {
		order, err := ordersRepository.SaveOrder(c.Request.Context(), orderRequest)
		if err != nil {
			switch err {
			case errors.InvalidMethod:
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case errors.InvalidOrder:
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			default:
				c.AbortWithError(http.StatusNotFound, err)
			}
			return
		}
		var wg sync.WaitGroup
		for _, topic := range [2]string{broker.UserStatusValidatorTopic, broker.UserBalanceValidatorTopic} {
			wg.Add(1)

			go func(t string, o entities.Order) {
				defer wg.Done()
				kafkaAdapter.PublishOrderMessageToTopic(t, o)
			}(topic, order)
		}

		wg.Wait()

		c.JSON(http.StatusOK, gin.H{"order_id": order.ID})
		return
	}
	c.String(http.StatusNotFound, "Something went wrong. Try again.")
}

func GetOrdersByUserId(c *gin.Context) {
	userId := c.Param("userId")
	offset, err := strconv.ParseInt(c.Query("offset"), 0, 64)
	if err != nil {
		offset = 0
	}
	limit, err := strconv.ParseInt(c.Query("limit"), 0, 64)
	if err != nil {
		limit = 100
	}

	orders, err := ordersRepository.OrdersDBAdapter.GetOrdersByUserId(c.Request.Context(), userId, offset, limit)

	var ordersResponse []dto.OrderResponse

	for _, order := range orders {
		ordersResponse = append(
			ordersResponse,
			dto.OrderResponse{
				Id:          order.ID,
				UserId:      order.UserId,
				Value:       order.Value,
				Installment: int64(order.Method.Installment),
				Status:      order.Status,
				Method:      order.Method.Name,
			},
		)
	}
	c.JSON(http.StatusOK, ordersResponse)
}
