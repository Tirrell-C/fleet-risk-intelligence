package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AppError represents application-specific errors with context
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
	Internal   error  `json:"-"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %s (internal: %s)", e.Code, e.Message, e.Internal.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Error constructors for common scenarios
func DatabaseError(operation string, err error) *AppError {
	return &AppError{
		Code:       "database_error",
		Message:    fmt.Sprintf("Database operation failed: %s", operation),
		HTTPStatus: http.StatusInternalServerError,
		Internal:   err,
		Context: map[string]interface{}{
			"operation": operation,
		},
	}
}

func ValidationError(field, message string) *AppError {
	return &AppError{
		Code:       "validation_error",
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Context: map[string]interface{}{
			"field": field,
		},
	}
}

func NotFoundError(resource string, id interface{}) *AppError {
	return &AppError{
		Code:       "not_found",
		Message:    fmt.Sprintf("%s not found", resource),
		HTTPStatus: http.StatusNotFound,
		Context: map[string]interface{}{
			"resource": resource,
			"id":       id,
		},
	}
}

func TelemetryIngestionError(vehicleID uint, err error) *AppError {
	return &AppError{
		Code:       "telemetry_ingestion_error",
		Message:    "Failed to ingest telemetry data",
		HTTPStatus: http.StatusInternalServerError,
		Internal:   err,
		Context: map[string]interface{}{
			"vehicle_id": vehicleID,
		},
	}
}

func RiskProcessingError(eventID uint, err error) *AppError {
	return &AppError{
		Code:       "risk_processing_error",
		Message:    "Failed to process risk event",
		HTTPStatus: http.StatusInternalServerError,
		Internal:   err,
		Context: map[string]interface{}{
			"event_id": eventID,
		},
	}
}

// ErrorHandler middleware for consistent error handling
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Check if it's our custom AppError
			if appErr, ok := err.Err.(*AppError); ok {
				// Log the error with context
				logEntry := logrus.WithFields(logrus.Fields{
					"error_code": appErr.Code,
					"context":    appErr.Context,
				})

				if appErr.Internal != nil {
					logEntry = logEntry.WithError(appErr.Internal)
				}

				logEntry.Error(appErr.Message)

				// Return structured error response
				c.JSON(appErr.HTTPStatus, gin.H{
					"error":   appErr.Code,
					"message": appErr.Message,
					"context": appErr.Context,
				})
				return
			}

			// Handle generic errors
			logrus.WithError(err.Err).Error("Unhandled error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "An internal error occurred",
			})
		}
	}
}

// LogAndAbort logs the error and aborts the request with proper error response
func LogAndAbort(c *gin.Context, err *AppError) {
	logEntry := logrus.WithFields(logrus.Fields{
		"error_code": err.Code,
		"context":    err.Context,
	})

	if err.Internal != nil {
		logEntry = logEntry.WithError(err.Internal)
	}

	logEntry.Error(err.Message)

	c.JSON(err.HTTPStatus, gin.H{
		"error":   err.Code,
		"message": err.Message,
		"context": err.Context,
	})
	c.Abort()
}

// WrapDatabaseError wraps database errors with additional context
func WrapDatabaseError(operation string, err error, context map[string]interface{}) *AppError {
	appErr := DatabaseError(operation, err)
	if context != nil {
		for k, v := range context {
			appErr.Context[k] = v
		}
	}
	return appErr
}