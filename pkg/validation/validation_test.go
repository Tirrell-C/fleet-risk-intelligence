package validation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateVehicleID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		vehicleID      string
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name:           "Valid vehicle ID",
			vehicleID:      "123",
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
		},
		{
			name:           "Empty vehicle ID",
			vehicleID:      " ",
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name:           "Invalid vehicle ID - not a number",
			vehicleID:      "abc",
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name:           "Invalid vehicle ID - zero",
			vehicleID:      "0",
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name:           "Invalid vehicle ID - negative",
			vehicleID:      "-1",
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			router.GET("/test/:vehicle_id", ValidateVehicleID(), func(c *gin.Context) {
				vehicleID, exists := c.Get("vehicle_id")
				if !exists {
					return // middleware should have aborted
				}
				assert.IsType(t, uint(0), vehicleID)
				c.JSON(http.StatusOK, gin.H{"vehicle_id": vehicleID})
			})

			req, _ := http.NewRequest("GET", "/test/"+tt.vehicleID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.shouldAbort {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "validation_failed", response["error"])
			}
		})
	}
}

func TestValidateTelemetryPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name: "Valid location payload",
			payload: map[string]interface{}{
				"vehicle_id":  1,
				"event_type":  "location",
				"timestamp":   time.Now().Format(time.RFC3339),
				"latitude":    37.7749,
				"longitude":   -122.4194,
				"speed":       55.0,
				"acceleration": 2.5,
			},
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
		},
		{
			name: "Valid speed payload without coordinates",
			payload: map[string]interface{}{
				"vehicle_id": 1,
				"event_type": "speed",
				"timestamp":  time.Now().Format(time.RFC3339),
				"speed":      65.0,
			},
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
		},
		{
			name: "Invalid - missing required fields",
			payload: map[string]interface{}{
				"vehicle_id": 1,
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "Invalid - zero vehicle_id",
			payload: map[string]interface{}{
				"vehicle_id": 0,
				"event_type": "location",
				"timestamp":  time.Now().Format(time.RFC3339),
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "Invalid - invalid event_type",
			payload: map[string]interface{}{
				"vehicle_id": 1,
				"event_type": "invalid_event",
				"timestamp":  time.Now().Format(time.RFC3339),
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "Invalid - future timestamp",
			payload: map[string]interface{}{
				"vehicle_id": 1,
				"event_type": "location",
				"timestamp":  time.Now().Add(10 * time.Minute).Format(time.RFC3339),
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "Invalid - latitude out of range",
			payload: map[string]interface{}{
				"vehicle_id": 1,
				"event_type": "location",
				"timestamp":  time.Now().Format(time.RFC3339),
				"latitude":   100.0,
				"longitude":  -122.4194,
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "Invalid - longitude out of range",
			payload: map[string]interface{}{
				"vehicle_id": 1,
				"event_type": "location",
				"timestamp":  time.Now().Format(time.RFC3339),
				"latitude":   37.7749,
				"longitude":  200.0,
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "Invalid - speed out of range",
			payload: map[string]interface{}{
				"vehicle_id": 1,
				"event_type": "speed",
				"timestamp":  time.Now().Format(time.RFC3339),
				"speed":      350.0,
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "Invalid - acceleration out of range",
			payload: map[string]interface{}{
				"vehicle_id":   1,
				"event_type":   "acceleration",
				"timestamp":    time.Now().Format(time.RFC3339),
				"acceleration": 25.0,
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			router.POST("/test", ValidateTelemetryPayload(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "valid"})
			})

			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.shouldAbort {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "validation_failed", response["error"])
			}
		})
	}
}

func TestRequireCoordinates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name: "Location event with coordinates",
			payload: map[string]interface{}{
				"event_type": "location",
				"latitude":   37.7749,
				"longitude":  -122.4194,
			},
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
		},
		{
			name: "Non-location event without coordinates",
			payload: map[string]interface{}{
				"event_type": "speed",
			},
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
		},
		{
			name: "Location event missing latitude",
			payload: map[string]interface{}{
				"event_type": "location",
				"longitude":  -122.4194,
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "Location event missing longitude",
			payload: map[string]interface{}{
				"event_type": "location",
				"latitude":   37.7749,
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "Location event missing both coordinates",
			payload: map[string]interface{}{
				"event_type": "location",
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			router.POST("/test", RequireCoordinates(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "valid"})
			})

			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.shouldAbort {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "validation_failed", response["error"])
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	errors := ValidationErrors{
		{Field: "vehicle_id", Message: "vehicle_id is required"},
		{Field: "event_type", Message: "invalid event type"},
	}

	assert.Equal(t, "validation failed: vehicle_id is required", errors.Error())

	// Test empty errors
	emptyErrors := ValidationErrors{}
	assert.Equal(t, "validation errors", emptyErrors.Error())
}