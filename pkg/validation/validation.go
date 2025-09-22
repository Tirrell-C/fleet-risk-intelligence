package validation

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation errors"
	}
	return fmt.Sprintf("validation failed: %s", ve[0].Message)
}

// ValidateVehicleID validates vehicle ID parameter
func ValidateVehicleID() gin.HandlerFunc {
	return func(c *gin.Context) {
		vehicleIDStr := strings.TrimSpace(c.Param("vehicle_id"))
		if vehicleIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_failed",
				"message": "vehicle_id is required",
			})
			c.Abort()
			return
		}

		vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 32)
		if err != nil || vehicleID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_failed",
				"message": "vehicle_id must be a valid positive integer",
			})
			c.Abort()
			return
		}

		// Store parsed value for handlers to use
		c.Set("vehicle_id", uint(vehicleID))
		c.Next()
	}
}

// ValidateTelemetryPayload validates telemetry payload data
func ValidateTelemetryPayload() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			VehicleID    uint      `json:"vehicle_id" binding:"required"`
			EventType    string    `json:"event_type" binding:"required"`
			Timestamp    time.Time `json:"timestamp" binding:"required"`
			Latitude     *float64  `json:"latitude"`
			Longitude    *float64  `json:"longitude"`
			Speed        *float64  `json:"speed"`
			Acceleration *float64  `json:"acceleration"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_failed",
				"message": "Invalid JSON payload: " + err.Error(),
			})
			c.Abort()
			return
		}

		errors := ValidationErrors{}

		// Validate vehicle ID
		if payload.VehicleID == 0 {
			errors = append(errors, ValidationError{
				Field:   "vehicle_id",
				Message: "vehicle_id must be greater than 0",
			})
		}

		// Validate event type
		validEventTypes := map[string]bool{
			"location":           true,
			"speed":              true,
			"acceleration":       true,
			"harsh_braking":      true,
			"engine_status":      true,
			"fuel_level":         true,
		}
		if !validEventTypes[payload.EventType] {
			errors = append(errors, ValidationError{
				Field:   "event_type",
				Message: "event_type must be one of: location, speed, acceleration, harsh_braking, engine_status, fuel_level",
			})
		}

		// Validate timestamp is not too far in the future
		if payload.Timestamp.After(time.Now().Add(5 * time.Minute)) {
			errors = append(errors, ValidationError{
				Field:   "timestamp",
				Message: "timestamp cannot be more than 5 minutes in the future",
			})
		}

		// Validate coordinates if provided
		if payload.Latitude != nil {
			if *payload.Latitude < -90 || *payload.Latitude > 90 {
				errors = append(errors, ValidationError{
					Field:   "latitude",
					Message: "latitude must be between -90 and 90",
				})
			}
		}

		if payload.Longitude != nil {
			if *payload.Longitude < -180 || *payload.Longitude > 180 {
				errors = append(errors, ValidationError{
					Field:   "longitude",
					Message: "longitude must be between -180 and 180",
				})
			}
		}

		// Validate speed if provided
		if payload.Speed != nil {
			if *payload.Speed < 0 || *payload.Speed > 300 { // 300 mph seems reasonable max
				errors = append(errors, ValidationError{
					Field:   "speed",
					Message: "speed must be between 0 and 300 mph",
				})
			}
		}

		// Validate acceleration if provided
		if payload.Acceleration != nil {
			if *payload.Acceleration < -20 || *payload.Acceleration > 20 { // Reasonable g-force limits
				errors = append(errors, ValidationError{
					Field:   "acceleration",
					Message: "acceleration must be between -20 and 20 m/s²",
				})
			}
		}

		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_failed",
				"message": "Validation failed",
				"errors":  errors,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireCoordinates ensures latitude and longitude are both provided for location events
func RequireCoordinates() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			EventType string   `json:"event_type"`
			Latitude  *float64 `json:"latitude"`
			Longitude *float64 `json:"longitude"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.Next() // Let other validation handle JSON errors
			return
		}

		if payload.EventType == "location" {
			if payload.Latitude == nil || payload.Longitude == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "validation_failed",
					"message": "Location events must include both latitude and longitude",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}