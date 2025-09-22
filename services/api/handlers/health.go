package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

// HealthCheck returns a health check handler
func HealthCheck(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Services:  make(map[string]string),
			Version:   "1.0.0",
		}

		// Check database connection
		sqlDB, err := db.DB()
		if err != nil {
			response.Status = "unhealthy"
			response.Services["database"] = "error: " + err.Error()
			c.JSON(http.StatusServiceUnavailable, response)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			response.Status = "unhealthy"
			response.Services["database"] = "error: " + err.Error()
			c.JSON(http.StatusServiceUnavailable, response)
			return
		}

		response.Services["database"] = "healthy"

		c.JSON(http.StatusOK, response)
	}
}