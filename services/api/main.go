package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/config"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/database"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/models"
	"github.com/Tirrell-C/fleet-risk-intelligence/services/api/graph"
	"github.com/Tirrell-C/fleet-risk-intelligence/services/api/handlers"
)

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

	// Run migrations
	if err := models.Migrate(db); err != nil {
		logrus.WithError(err).Fatal("Failed to run database migrations")
	}

	// Setup GraphQL server
	resolver := &graph.Resolver{
		DB:     db,
		Config: cfg,
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	// Add WebSocket support for subscriptions
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
	})

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

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck(db))

	// GraphQL endpoints
	router.POST("/graphql", gin.WrapH(srv))
	router.GET("/graphql", gin.WrapH(srv))

	// GraphQL playground (development only)
	if cfg.Server.Env == "development" {
		router.GET("/playground", gin.WrapH(playground.Handler("GraphQL Playground", "/graphql")))
		logrus.Info("GraphQL Playground available at http://localhost:" + cfg.Server.Port + "/playground")
	}

	// Setup HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logrus.WithField("port", cfg.Server.Port).Info("Starting GraphQL API server")
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