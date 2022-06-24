package handlers

import (
	"net/http"

	"github.com/evertontomalok/distributed-system-go/internal/infra/database/postgres"
	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func Readiness(c *gin.Context) {
	if !postgres.Started {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	c.String(http.StatusOK, "ok")
}
