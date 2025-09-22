package main

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/auth"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/server"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/models"
	"github.com/Tirrell-C/fleet-risk-intelligence/services/api/graph"
)

func main() {
	// Initialize base server with common setup
	baseServer, err := server.NewBaseServer("api")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize server")
	}

	// Initialize JWT manager
	jwtSecret := baseServer.Config.Server.JWTSecret
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
		logrus.Warn("Using default JWT secret - change this in production!")
	}
	jwtManager := auth.NewJWTManager(jwtSecret, 24*time.Hour)
	authMiddleware := auth.NewAuthMiddleware(jwtManager)

	// Add GraphQL endpoint
	setupGraphQL(baseServer, authMiddleware)

	// Add basic REST endpoints
	setupRoutes(baseServer, authMiddleware)

	// Start server
	if err := baseServer.Start(baseServer.Config.Server.Port); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}

	// Wait for shutdown
	baseServer.WaitForShutdown()
}

func setupGraphQL(server *server.BaseServer, authMiddleware *auth.AuthMiddleware) {
	// Create GraphQL resolver with database access
	resolver := &graph.Resolver{
		DB:     server.DB,
		Config: server.Config,
	}

	// Create GraphQL handler
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Add GraphQL endpoint with authentication
	server.Router.POST("/graphql", authMiddleware.RequireAuth(), func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	// Add GraphQL playground for development
	if server.Config.Server.Env == "development" {
		server.Router.GET("/playground", func(c *gin.Context) {
			playground.Handler("GraphQL playground", "/graphql").ServeHTTP(c.Writer, c.Request)
		})
		logrus.Info("GraphQL playground available at /playground")
	}

	logrus.Info("GraphQL endpoint available at /graphql")
}

func setupRoutes(server *server.BaseServer, authMiddleware *auth.AuthMiddleware) {
	api := server.Router.Group("/api/v1")
	api.Use(authMiddleware.RequireAuth())

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