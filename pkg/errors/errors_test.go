package errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAppError(t *testing.T) {
	// Test without internal error
	err := &AppError{
		Code:       "test_error",
		Message:    "Test error message",
		HTTPStatus: http.StatusBadRequest,
	}

	assert.Equal(t, "test_error: Test error message", err.Error())

	// Test with internal error
	internalErr := errors.New("internal error")
	err = &AppError{
		Code:       "test_error",
		Message:    "Test error message",
		HTTPStatus: http.StatusBadRequest,
		Internal:   internalErr,
	}

	assert.Equal(t, "test_error: Test error message (internal: internal error)", err.Error())
}

func TestDatabaseError(t *testing.T) {
	internalErr := errors.New("connection failed")
	err := DatabaseError("create_user", internalErr)

	assert.Equal(t, "database_error", err.Code)
	assert.Equal(t, "Database operation failed: create_user", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus)
	assert.Equal(t, internalErr, err.Internal)
	assert.Equal(t, "create_user", err.Context["operation"])
}

func TestValidationError(t *testing.T) {
	err := ValidationError("email", "invalid email format")

	assert.Equal(t, "validation_error", err.Code)
	assert.Equal(t, "invalid email format", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.HTTPStatus)
	assert.Equal(t, "email", err.Context["field"])
}

func TestNotFoundError(t *testing.T) {
	err := NotFoundError("user", 123)

	assert.Equal(t, "not_found", err.Code)
	assert.Equal(t, "user not found", err.Message)
	assert.Equal(t, http.StatusNotFound, err.HTTPStatus)
	assert.Equal(t, "user", err.Context["resource"])
	assert.Equal(t, 123, err.Context["id"])
}

func TestTelemetryIngestionError(t *testing.T) {
	internalErr := errors.New("database insert failed")
	err := TelemetryIngestionError(456, internalErr)

	assert.Equal(t, "telemetry_ingestion_error", err.Code)
	assert.Equal(t, "Failed to ingest telemetry data", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus)
	assert.Equal(t, internalErr, err.Internal)
	assert.Equal(t, uint(456), err.Context["vehicle_id"])
}

func TestRiskProcessingError(t *testing.T) {
	internalErr := errors.New("processing failed")
	err := RiskProcessingError(789, internalErr)

	assert.Equal(t, "risk_processing_error", err.Code)
	assert.Equal(t, "Failed to process risk event", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus)
	assert.Equal(t, internalErr, err.Internal)
	assert.Equal(t, uint(789), err.Context["event_id"])
}

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle AppError", func(t *testing.T) {
		router := gin.New()
		router.Use(ErrorHandler())

		router.GET("/test", func(c *gin.Context) {
			err := ValidationError("test_field", "test validation error")
			c.Error(err)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "validation_error")
		assert.Contains(t, w.Body.String(), "test validation error")
	})

	t.Run("Handle generic error", func(t *testing.T) {
		router := gin.New()
		router.Use(ErrorHandler())

		router.GET("/test", func(c *gin.Context) {
			c.Error(errors.New("generic error"))
			c.JSON(http.StatusInternalServerError, gin.H{"should": "not reach here"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal_error")
		assert.Contains(t, w.Body.String(), "An internal error occurred")
	})

	t.Run("No error", func(t *testing.T) {
		router := gin.New()
		router.Use(ErrorHandler())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})
}

func TestLogAndAbort(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	router.GET("/test", func(c *gin.Context) {
		err := NotFoundError("resource", "test_id")
		LogAndAbort(c, err)
		c.JSON(http.StatusOK, gin.H{"should": "not reach here"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "not_found")
	assert.Contains(t, w.Body.String(), "resource not found")
}

func TestWrapDatabaseError(t *testing.T) {
	internalErr := errors.New("connection timeout")
	context := map[string]interface{}{
		"table":  "users",
		"action": "insert",
	}

	err := WrapDatabaseError("create_user", internalErr, context)

	assert.Equal(t, "database_error", err.Code)
	assert.Equal(t, "Database operation failed: create_user", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus)
	assert.Equal(t, internalErr, err.Internal)
	assert.Equal(t, "create_user", err.Context["operation"])
	assert.Equal(t, "users", err.Context["table"])
	assert.Equal(t, "insert", err.Context["action"])

	// Test with nil context
	err2 := WrapDatabaseError("delete_user", internalErr, nil)
	assert.Equal(t, "create_user", err.Context["operation"]) // Original context preserved
	assert.Equal(t, "delete_user", err2.Context["operation"]) // New error has correct operation
}