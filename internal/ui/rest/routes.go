package rest

import (
	"net/http"

	"github.com/evertontomalok/distributed-system-go/internal/ui/rest/handlers"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

var healthCheck = []Route{
	{
		"/health",
		http.MethodGet,
		handlers.Health,
	},
	{
		"/readiness",
		http.MethodGet,
		handlers.Readiness,
	},
}

var routes = []Route{
	{
		"/orders",
		http.MethodPost,
		handlers.PostOrder,
	},
	{
		"/orders/:userId",
		http.MethodGet,
		handlers.GetOrdersByUserId,
	},
	{
		"/order/:orderId",
		http.MethodGet,
		handlers.GetOrderById,
	},
	{
		"/orders/user/test",
		http.MethodDelete,
		handlers.DeleteOrdersFromUserTest,
	},
	{
		"/methods",
		http.MethodGet,
		handlers.GetAllMethods,
	},
}
