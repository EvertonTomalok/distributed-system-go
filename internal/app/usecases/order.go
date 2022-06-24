package usecases

import (
	"net/http"
	"strconv"

	domainErrors "github.com/evertontomalok/distributed-system-go/internal/core/domain/errors"
	"github.com/evertontomalok/distributed-system-go/internal/core/dto"
	ordersRepository "github.com/evertontomalok/distributed-system-go/internal/core/orders"
	"github.com/evertontomalok/distributed-system-go/internal/infra/services/aws"
	"github.com/evertontomalok/distributed-system-go/pkg/helpers"
	"github.com/evertontomalok/distributed-system-go/pkg/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func CreateOrder(c *gin.Context) *gin.Context {
	if status := utils.CheckFeatureFlag(dto.PostOrderFlag, c); !status {
		log.Info("Post Order Flag is disabled.")
		c.AbortWithStatus(http.StatusBadGateway)
		return c
	}
	orderRequest := dto.OrderRequest{}

	err := c.ShouldBind(&orderRequest)

	if err == nil {
		order, err := ordersRepository.SaveOrder(c.Request.Context(), orderRequest)
		if err != nil {
			switch err {
			case domainErrors.InvalidMethod, domainErrors.InvalidOrder:
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				if err := c.AbortWithError(http.StatusNotFound, err); err != nil {
					log.Fatalf("Something went wrong: %+v", err)
				}
			}
			return c
		}

		if err := helpers.TriggerValidation(order); err != nil {
			aws.SendErrorToCloudWatch(c, err)
			c.String(http.StatusInternalServerError, "Something went wrong. Try again.")
			return c
		}

		orderResponse := dto.OrderResponse{
			Id:          order.ID,
			Status:      order.Status,
			Value:       order.Value,
			UserId:      orderRequest.UserId,
			Installment: orderRequest.Installment,
			Method:      orderRequest.Method,
		}
		c.JSON(http.StatusCreated, orderResponse)
		return c
	}
	c.String(http.StatusUnprocessableEntity, "Something went wrong. Try again.")
	return c
}

func GetAllOrdersFromUserById(c *gin.Context) *gin.Context {
	if status := utils.CheckFeatureFlag(dto.GetOrdersFromUserFlag, c); !status {
		log.Info("Get Orders from User flag is disabled.")
		c.AbortWithStatus(http.StatusBadGateway)
		return c
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

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return c
	}

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
	return c
}

func GetOrderById(c *gin.Context) *gin.Context {
	if status := utils.CheckFeatureFlag(dto.GetOrderByIdFlag, c); !status {
		log.Info("Get Order by id flag is disabled.")
		c.AbortWithStatus(http.StatusBadGateway)
		return c
	}
	orderId := c.Param("orderId")

	order, err := ordersRepository.OrdersDBAdapter.GetOrderById(c.Request.Context(), orderId)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return c
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
	return c
}
