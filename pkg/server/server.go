package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/config"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/database"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/models"
)

// BaseServer provides common server functionality
type BaseServer struct {
	DB     *gorm.DB
	Config *config.Config
	Router *gin.Engine
	server *http.Server
}

// NewBaseServer creates a new base server with common setup
func NewBaseServer(serviceName string) (*BaseServer, error) {
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
		return nil, err
	}

	// Run migrations
	if err := models.Migrate(db); err != nil {
		return nil, err
	}

	// Setup Gin router with common middleware
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// CORS middleware with environment-based configuration
	router.Use(corsMiddleware(cfg.Server.Env))

	// Create base server
	bs := &BaseServer{
		DB:     db,
		Config: cfg,
		Router: router,
	}

	// Add common health check
	router.GET("/health", bs.healthCheck(serviceName))

	return bs, nil
}

// Start starts the server on the specified port
func (bs *BaseServer) Start(port string) error {
	bs.server = &http.Server{
		Addr:    ":" + port,
		Handler: bs.Router,
	}

	// Start server in a goroutine
	go func() {
		logrus.WithField("port", port).Info("Starting server")
		if err := bs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("Failed to start server")
		}
	}()

	return nil
}

// WaitForShutdown waits for interrupt signal and gracefully shuts down
func (bs *BaseServer) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := bs.server.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("Server forced to shutdown")
	}

	logrus.Info("Server exited")
}

// corsMiddleware returns CORS middleware based on environment
func corsMiddleware(env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := "*"
		if env == "production" {
			// In production, specify allowed origins
			origin = "https://yourdomain.com"
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// healthCheck returns a health check handler
func (bs *BaseServer) healthCheck(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Test database connection
		sqlDB, err := bs.DB.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "unhealthy",
				"service":   serviceName,
				"error":     "database connection failed",
				"timestamp": time.Now(),
			})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "unhealthy",
				"service":   serviceName,
				"error":     "database ping failed",
				"timestamp": time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   serviceName,
			"timestamp": time.Now(),
		})
	}
}

// setupLogging configures logging based on environment
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