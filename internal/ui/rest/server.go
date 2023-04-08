package rest

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/evertontomalok/distributed-system-go/docs"

	config "github.com/evertontomalok/distributed-system-go/internal/app"
	"github.com/evertontomalok/distributed-system-go/pkg/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	router := gin.Default()

	injectRoutes(router)

	return router
}

// @title           Go Distributed System
// @version         1.0
// @description     This is a sample server Event Driven System.

// @contact.name   Everton Tomalok
// @contact.email  evertontomalok123@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api
func RunServer(ctx context.Context, config config.Config) {
	done := utils.MakeDoneSignal()

	server := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: Router(),
	}

	go func() {
		log.Printf("Server started at %s:%s", config.Host, config.Port)

		if err := server.ListenAndServe(); err != nil {
			log.Panicf("Error trying to start server. %+v", err)
		}
	}()

	<-done
	log.Println("Stopping server...")
}

func injectRoutes(router *gin.Engine) {
	docs.SwaggerInfo.Title = "Go Distributed System"
	docs.SwaggerInfo.Description = "Go Distributed System."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	for _, route := range healthCheck {
		router.Handle(route.Method, route.Path, route.Handler)
	}

	apiGroup := router.Group("/api")
	for _, route := range routes {
		apiGroup.Handle(route.Method, route.Path, route.Handler)
	}
}
