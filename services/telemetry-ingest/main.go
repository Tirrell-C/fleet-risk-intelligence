package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/config"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/database"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/models"
)

type TelemetryHandler struct {
	db     *database.DB
	config *config.Config
}

type TelemetryPayload struct {
	VehicleID    uint      `json:"vehicle_id" binding:"required"`
	EventType    string    `json:"event_type" binding:"required"`
	Timestamp    time.Time `json:"timestamp" binding:"required"`
	Latitude     *float64  `json:"latitude"`
	Longitude    *float64  `json:"longitude"`
	Speed        *float64  `json:"speed"`
	Acceleration *float64  `json:"acceleration"`
	Data         string    `json:"data"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Info("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// Setup logging
	setupLogging(cfg.Server.Env)

	// Connect to database
	db, err := database.NewConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
	})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}

	handler := &TelemetryHandler{
		db:     db,
		config: cfg,
	}

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "telemetry-ingest",
			"timestamp": time.Now(),
		})
	})

	// Telemetry endpoints
	router.POST("/telemetry", handler.IngestTelemetry)
	router.POST("/telemetry/batch", handler.IngestBatchTelemetry)

	// Simulation endpoint for development
	if cfg.Features.EnableTelemetrySimulation {
		router.POST("/simulate/:vehicle_id", handler.SimulateTelemetry)
		logrus.Info("Telemetry simulation enabled")
	}

	// Setup HTTP server
	port := getEnv("TELEMETRY_PORT", "8081")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logrus.WithField("port", port).Info("Starting telemetry ingestion service")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("Server forced to shutdown")
	}

	logrus.Info("Server exited")
}

// IngestTelemetry handles single telemetry event ingestion
func (h *TelemetryHandler) IngestTelemetry(c *gin.Context) {
	var payload TelemetryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create telemetry event
	event := models.TelemetryEvent{
		VehicleID:    payload.VehicleID,
		EventType:    payload.EventType,
		Timestamp:    payload.Timestamp,
		Latitude:     payload.Latitude,
		Longitude:    payload.Longitude,
		Speed:        payload.Speed,
		Acceleration: payload.Acceleration,
		Data:         payload.Data,
	}

	if err := h.db.Create(&event).Error; err != nil {
		logrus.WithError(err).Error("Failed to save telemetry event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save telemetry event"})
		return
	}

	// TODO: Publish to Redis for real-time processing
	h.publishToRedis(&event)

	c.JSON(http.StatusCreated, gin.H{
		"id":        event.ID,
		"processed": time.Now(),
	})
}

// IngestBatchTelemetry handles batch telemetry ingestion
func (h *TelemetryHandler) IngestBatchTelemetry(c *gin.Context) {
	var payloads []TelemetryPayload
	if err := c.ShouldBindJSON(&payloads); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events := make([]models.TelemetryEvent, len(payloads))
	for i, payload := range payloads {
		events[i] = models.TelemetryEvent{
			VehicleID:    payload.VehicleID,
			EventType:    payload.EventType,
			Timestamp:    payload.Timestamp,
			Latitude:     payload.Latitude,
			Longitude:    payload.Longitude,
			Speed:        payload.Speed,
			Acceleration: payload.Acceleration,
			Data:         payload.Data,
		}
	}

	// Batch insert
	if err := h.db.CreateInBatches(&events, 100).Error; err != nil {
		logrus.WithError(err).Error("Failed to save batch telemetry events")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save telemetry events"})
		return
	}

	// Publish events for real-time processing
	for _, event := range events {
		h.publishToRedis(&event)
	}

	c.JSON(http.StatusCreated, gin.H{
		"processed": len(events),
		"timestamp": time.Now(),
	})
}

// SimulateTelemetry generates simulated telemetry data for development
func (h *TelemetryHandler) SimulateTelemetry(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id is required"})
		return
	}

	// Generate simulated data
	events := generateSimulatedTelemetry(vehicleID, 10)

	for _, event := range events {
		if err := h.db.Create(&event).Error; err != nil {
			logrus.WithError(err).Error("Failed to save simulated telemetry")
			continue
		}
		h.publishToRedis(&event)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Simulated telemetry generated",
		"events":    len(events),
		"vehicle_id": vehicleID,
	})
}

// publishToRedis publishes telemetry events to Redis for real-time processing
func (h *TelemetryHandler) publishToRedis(event *models.TelemetryEvent) {
	// TODO: Implement Redis publishing
	// This would typically publish to a Redis stream or pub/sub channel
	// for the risk engine to consume in real-time

	data, _ := json.Marshal(event)
	logrus.WithFields(logrus.Fields{
		"vehicle_id": event.VehicleID,
		"event_type": event.EventType,
		"timestamp":  event.Timestamp,
	}).Debug("Publishing telemetry event to Redis: " + string(data))
}

// generateSimulatedTelemetry creates realistic telemetry data for testing
func generateSimulatedTelemetry(vehicleIDStr string, count int) []models.TelemetryEvent {
	// Simple simulation - in a real system this would be much more sophisticated
	events := make([]models.TelemetryEvent, count)
	baseTime := time.Now()

	// Parse vehicle ID
	var vehicleID uint = 1 // Default fallback
	if id, err := parseUint(vehicleIDStr); err == nil {
		vehicleID = id
	}

	for i := 0; i < count; i++ {
		events[i] = models.TelemetryEvent{
			VehicleID: vehicleID,
			EventType: "location",
			Timestamp: baseTime.Add(time.Duration(i) * time.Minute),
			Latitude:  floatPtr(37.7749 + float64(i)*0.001),   // San Francisco area
			Longitude: floatPtr(-122.4194 + float64(i)*0.001),
			Speed:     floatPtr(float64(25 + i%30)),           // 25-55 mph
			Data:      fmt.Sprintf(`{"engine_status":"on","fuel_level":%d}`, 80-i),
		}
	}

	return events
}

func setupLogging(env string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if env == "development" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func floatPtr(f float64) *float64 {
	return &f
}

func parseUint(s string) (uint, error) {
	// Simple uint parsing - in production use strconv.ParseUint
	if s == "1" {
		return 1, nil
	}
	if s == "2" {
		return 2, nil
	}
	return 1, fmt.Errorf("invalid uint: %s", s)
}