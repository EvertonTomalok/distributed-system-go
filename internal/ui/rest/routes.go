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
		"/",
		http.MethodGet,
		handlers.Home,
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
		"/methods",
		http.MethodGet,
		handlers.GetAllMethods,
	},
}
