package rest

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	config "github.com/evertontomalok/distributed-system-go/internal/app"
	"github.com/evertontomalok/distributed-system-go/internal/ui/rest/handlers"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.GET("/", handlers.Home)
	router.POST("/orders", handlers.PostOrder)
	router.GET("/orders/:userId", handlers.GetOrdersByUserId)
	router.GET("/methods", handlers.GetAllMethods)

	return router
}

func RunServer(ctx context.Context, config config.Config) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

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
