package rest

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/evertontomalok/distributed-system-go/src/config"
	"github.com/evertontomalok/distributed-system-go/src/controllers"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.GET("/", controllers.Home)
	router.GET("/orders/:userId", controllers.GetOrdersByUserId)

	return router
}

func RunServer(ctx context.Context, serverConfig config.ServerConfig) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    net.JoinHostPort(serverConfig.Host, serverConfig.Port),
		Handler: Router(),
	}

	go func() {
		log.Printf("Server started at %s:%s", serverConfig.Host, serverConfig.Port)

		if err := server.ListenAndServe(); err != nil {
			log.Panicf("Error trying to start server. %+v", err)
		}
	}()

	<-done
	log.Println("Stopping server...")
}
