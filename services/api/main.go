package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/server"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/models"
)

func main() {
	// Initialize base server with common setup
	baseServer, err := server.NewBaseServer("api")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize server")
	}

	// Add basic REST endpoints for now (GraphQL can be added later)
	setupRoutes(baseServer)

	// Start server
	if err := baseServer.Start(baseServer.Config.Server.Port); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}

	// Wait for shutdown
	baseServer.WaitForShutdown()
}

func setupRoutes(server *server.BaseServer) {
	api := server.Router.Group("/api/v1")

	// Vehicle endpoints
	api.GET("/vehicles", getVehicles(server))
	api.GET("/vehicles/:id", getVehicle(server))

	// Fleet endpoints
	api.GET("/fleets", getFleets(server))
	api.GET("/fleets/:id", getFleet(server))

	// Driver endpoints
	api.GET("/drivers", getDrivers(server))
	api.GET("/drivers/:id", getDriver(server))

	// Risk events
	api.GET("/risk-events", getRiskEvents(server))
	api.GET("/vehicles/:id/risk-events", getVehicleRiskEvents(server))

	// Alerts
	api.GET("/alerts", getAlerts(server))
}

func getVehicles(server *server.BaseServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var vehicles []models.Vehicle
		if err := server.DB.Preload("Fleet").Preload("Driver").Find(&vehicles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicles"})
			return
		}
		c.JSON(http.StatusOK, vehicles)
	}
}

func getVehicle(server *server.BaseServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var vehicle models.Vehicle
		if err := server.DB.Preload("Fleet").Preload("Driver").First(&vehicle, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}
		c.JSON(http.StatusOK, vehicle)
	}
}

func getFleets(server *server.BaseServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var fleets []models.Fleet
		if err := server.DB.Find(&fleets).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fleets"})
			return
		}
		c.JSON(http.StatusOK, fleets)
	}
}

func getFleet(server *server.BaseServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var fleet models.Fleet
		if err := server.DB.First(&fleet, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Fleet not found"})
			return
		}
		c.JSON(http.StatusOK, fleet)
	}
}

func getDrivers(server *server.BaseServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var drivers []models.Driver
		if err := server.DB.Preload("Fleet").Find(&drivers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch drivers"})
			return
		}
		c.JSON(http.StatusOK, drivers)
	}
}

func getDriver(server *server.BaseServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var driver models.Driver
		if err := server.DB.Preload("Fleet").First(&driver, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
			return
		}
		c.JSON(http.StatusOK, driver)
	}
}

func getRiskEvents(server *server.BaseServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var events []models.RiskEvent
		if err := server.DB.Preload("Vehicle").Preload("Driver").Order("created_at desc").Limit(100).Find(&events).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch risk events"})
			return
		}
		c.JSON(http.StatusOK, events)
	}
}

func getVehicleRiskEvents(server *server.BaseServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		vehicleID := c.Param("id")
		var events []models.RiskEvent
		if err := server.DB.Preload("Vehicle").Preload("Driver").Where("vehicle_id = ?", vehicleID).Order("created_at desc").Limit(100).Find(&events).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicle risk events"})
			return
		}
		c.JSON(http.StatusOK, events)
	}
}

func getAlerts(server *server.BaseServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var alerts []models.Alert
		if err := server.DB.Preload("Fleet").Preload("Vehicle").Preload("Driver").Order("created_at desc").Limit(100).Find(&alerts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch alerts"})
			return
		}
		c.JSON(http.StatusOK, alerts)
	}
}