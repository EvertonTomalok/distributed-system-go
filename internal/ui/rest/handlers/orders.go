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

// API Create Order godoc
// @Summary Create Order
// @Description Create order
// @Tags order
// @Router /orders [post]
// @Param Order body dto.OrderRequest true "Order to create"
// @Produce json
// @Success 201 "{'order_id': 'someid'}"
// @Failure 400 "{'error': 'error description'}"
// @Failure 404 "Something went wrong. Try again."
func PostOrder(c *gin.Context) {
	orderRequest := dto.OrderRequest{}

	err := c.ShouldBind(&orderRequest)

	if err == nil {
		order, err := ordersRepository.SaveOrder(c.Request.Context(), orderRequest)
		if err != nil {
			switch err {
			case errors.InvalidMethod:
			case errors.InvalidOrder:
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.AbortWithError(http.StatusNotFound, err)
			}
			return
		}

		triggerValidation(order)

		c.JSON(http.StatusCreated, gin.H{"order_id": order.ID})
		return
	}
	c.String(http.StatusNotFound, "Something went wrong. Try again.")
}

// API Create Order godoc
// @Summary Get Orders
// @Description Get orders from User
// @Tags order
// @Router /orders/{userId} [get]
// @Param userId path string false "The user id to search"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Accept json
// @Produce json
// @Success 200 {object} []dto.OrderResponse
// @Failure 500 "Something went wrong"
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
	if limit > 100 {
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

func triggerValidation(order entities.Order) {
	var wg sync.WaitGroup
	for _, topic := range [2]string{broker.UserStatusValidatorTopic, broker.UserBalanceValidatorTopic} {
		wg.Add(1)

		go func(t string, o entities.Order) {
			defer wg.Done()
			kafkaAdapter.PublishOrderMessageToTopic(t, o)
		}(topic, order)
	}

	wg.Wait()
}
