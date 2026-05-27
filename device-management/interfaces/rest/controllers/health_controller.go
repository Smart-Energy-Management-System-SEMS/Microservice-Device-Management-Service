package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "UP",
		"service":   "device-management-service",
		"timestamp": time.Now().UTC(),
	})
}
