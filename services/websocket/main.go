package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/server"
)

func main() {
	// Initialize base server with common setup
	baseServer, err := server.NewBaseServer("websocket")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize server")
	}

	// Simple websocket endpoint placeholder
	baseServer.Router.GET("/ws", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "WebSocket service placeholder - real-time features would be implemented here",
			"service": "websocket",
		})
	})

	// Start server
	port := getEnv("WEBSOCKET_PORT", "8083")
	if err := baseServer.Start(port); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}

	// Wait for shutdown
	baseServer.WaitForShutdown()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}