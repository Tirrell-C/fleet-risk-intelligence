package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/config"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/errors"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/models"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/server"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/validation"
)

type TelemetryHandler struct {
	db     *gorm.DB
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
	// Initialize base server with common setup
	baseServer, err := server.NewBaseServer("telemetry-ingest")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize server")
	}

	handler := &TelemetryHandler{
		db:     baseServer.DB,
		config: baseServer.Config,
	}

	// Add error handling middleware
	baseServer.Router.Use(errors.ErrorHandler())

	// Telemetry endpoints with validation
	baseServer.Router.POST("/telemetry",
		validation.ValidateTelemetryPayload(),
		validation.RequireCoordinates(),
		handler.IngestTelemetry)
	baseServer.Router.POST("/telemetry/batch", handler.IngestBatchTelemetry)

	// Simulation endpoint for development
	if baseServer.Config.Features.EnableTelemetrySimulation {
		baseServer.Router.POST("/simulate/:vehicle_id",
			validation.ValidateVehicleID(),
			handler.SimulateTelemetry)
		logrus.Info("Telemetry simulation enabled")
	}

	// Start server
	port := getEnv("TELEMETRY_PORT", "8081")
	if err := baseServer.Start(port); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}

	// Wait for shutdown
	baseServer.WaitForShutdown()
}

// IngestTelemetry handles single telemetry event ingestion
func (h *TelemetryHandler) IngestTelemetry(c *gin.Context) {
	var payload TelemetryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		errors.LogAndAbort(c, errors.ValidationError("json_payload", "Invalid JSON payload: "+err.Error()))
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
		errors.LogAndAbort(c, errors.TelemetryIngestionError(payload.VehicleID, err))
		return
	}

	// Note: Real-time processing would be added here in production

	c.JSON(http.StatusCreated, gin.H{
		"id":        event.ID,
		"processed": time.Now(),
	})
}

// IngestBatchTelemetry handles batch telemetry ingestion
func (h *TelemetryHandler) IngestBatchTelemetry(c *gin.Context) {
	var payloads []TelemetryPayload
	if err := c.ShouldBindJSON(&payloads); err != nil {
		errors.LogAndAbort(c, errors.ValidationError("json_payload", "Invalid JSON batch payload: "+err.Error()))
		return
	}

	if len(payloads) == 0 {
		errors.LogAndAbort(c, errors.ValidationError("batch_size", "Batch cannot be empty"))
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
		errors.LogAndAbort(c, errors.WrapDatabaseError("batch_telemetry_insert", err, map[string]interface{}{
			"batch_size": len(events),
		}))
		return
	}

	// Note: Real-time processing would be added here in production

	c.JSON(http.StatusCreated, gin.H{
		"processed": len(events),
		"timestamp": time.Now(),
	})
}

// SimulateTelemetry generates simulated telemetry data for development
func (h *TelemetryHandler) SimulateTelemetry(c *gin.Context) {
	// Get validated vehicle ID from middleware
	vehicleID, exists := c.Get("vehicle_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "vehicle_id validation failed"})
		return
	}

	// Generate simulated data
	events := generateSimulatedTelemetry(vehicleID.(uint), 10)

	successCount := 0
	for _, event := range events {
		if err := h.db.Create(&event).Error; err != nil {
			logrus.WithError(err).WithField("vehicle_id", vehicleID).Error("Failed to save simulated telemetry event")
			continue
		}
		successCount++
		// Note: Real-time processing would be added here in production
	}

	if successCount == 0 {
		errors.LogAndAbort(c, errors.TelemetryIngestionError(vehicleID.(uint), fmt.Errorf("no events were successfully saved")))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Simulated telemetry generated",
		"events_created":  successCount,
		"total_attempted": len(events),
		"vehicle_id":      vehicleID,
	})
}


// generateSimulatedTelemetry creates realistic telemetry data for testing
func generateSimulatedTelemetry(vehicleID uint, count int) []models.TelemetryEvent {
	// Simple simulation - in a real system this would be much more sophisticated
	events := make([]models.TelemetryEvent, count)
	baseTime := time.Now()

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
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid uint: %s", s)
	}
	return uint(val), nil
}