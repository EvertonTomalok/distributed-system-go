package utils

import (
	"net/http"

	"github.com/evertontomalok/distributed-system-go/internal/app/shared"
	"github.com/gin-gonic/gin"
)

func CheckFeatureFlag(flagName string, c *gin.Context) bool {
	flagStatus := shared.Flags[flagName]

	if flagStatus == false {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return false
	}
	return true
}
