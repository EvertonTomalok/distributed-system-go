package handlers

import (
	"github.com/evertontomalok/distributed-system-go/internal/app/usecases"
	"github.com/gin-gonic/gin"
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
	usecases.CreateOrder(c)
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
	usecases.GetAllOrdersFromUserById(c)
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
	usecases.GetOrderById(c)
}
