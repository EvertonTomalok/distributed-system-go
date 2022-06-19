package handlers

import (
	"net/http"
	"strconv"

	"github.com/evertontomalok/distributed-system-go/internal/app/utils"
	"github.com/evertontomalok/distributed-system-go/internal/domain/broker"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	domainErrors "github.com/evertontomalok/distributed-system-go/internal/domain/core/errors"
	ordersRepository "github.com/evertontomalok/distributed-system-go/internal/domain/orders"
	kafkaAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// API Create Order godoc
// @Summary Create Order
// @Description Create order
// @Tags order
// @Router /orders [post]
// @Param Order body dto.OrderRequest true "Order to create"
// @Produce json
// @Success 201 {object} dto.OrderResponse
// @Failure 400 "{'error': 'error description'}"
// @Failure 404 "Something went wrong. Try again."
// @Failure 503 "Feature Flag is disabled."
func PostOrder(c *gin.Context) {
	if status := utils.CheckFeatureFlag(dto.PostOrderFlag, c); status == false {
		log.Info("Post Order Flag is disabled.")
		return
	}
	orderRequest := dto.OrderRequest{}

	err := c.ShouldBind(&orderRequest)

	if err == nil {
		order, err := ordersRepository.SaveOrder(c.Request.Context(), orderRequest)
		if err != nil {
			log.Infoln("err", err)
			switch err {
			case domainErrors.InvalidMethod, domainErrors.InvalidOrder:
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.AbortWithError(http.StatusNotFound, err)
			}
			return
		}

		triggerValidation(order)
		orderResponse := dto.OrderResponse{
			Id:          order.ID,
			Status:      order.Status,
			Value:       order.Value,
			UserId:      orderRequest.UserId,
			Installment: orderRequest.Installment,
			Method:      orderRequest.Method,
		}
		c.JSON(http.StatusCreated, orderResponse)
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
// @Success 200 {object} dto.OrdersResponse
// @Failure 500 "Something went wrong"
func GetOrdersByUserId(c *gin.Context) {
	if status := utils.CheckFeatureFlag(dto.GetOrdersFromUserFlag, c); status == false {
		log.Info("Get Orders from User flag is disabled.")
		return
	}

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

	var ordersArray []dto.OrderResponse = make([]dto.OrderResponse, 0)

	for _, order := range orders {
		ordersArray = append(
			ordersArray,
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
	ordersResponse := dto.OrdersResponse{
		Offset: offset,
		Limit:  limit,
		Total:  len(ordersArray),
		Orders: ordersArray,
	}
	c.JSON(http.StatusOK, ordersResponse)
}

// API Create Order godoc
// @Summary Get Order by id
// @Description Get order using its id
// @Tags order
// @Router /orders/{orderId} [get]
// @Param orderId path string false "The order id to search"
// @Accept json
// @Produce json
// @Success 200 {object} dto.OrderResponse
// @Failure 404 "Order not found"
// @Failure 500 "Something went wrong"
func GetOrderById(c *gin.Context) {
	if status := utils.CheckFeatureFlag(dto.GetOrderByIdFlag, c); status == false {
		log.Info("Get Order by id flag is disabled.")
		return
	}
	orderId := c.Param("orderId")

	order, err := ordersRepository.OrdersDBAdapter.GetOrderById(c.Request.Context(), orderId)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	orderResponse := dto.OrderResponse{
		Id:          order.ID,
		UserId:      order.UserId,
		Status:      order.Status,
		Value:       order.Value,
		Method:      order.Method.Name,
		Installment: order.MethodId,
	}

	c.JSON(http.StatusOK, orderResponse)
}

func triggerValidation(order entities.Order) {
	kafkaAdapter.PublishOrderMessageToTopic(broker.OrchestratorTopic, order, dto.StartEvent)
}
