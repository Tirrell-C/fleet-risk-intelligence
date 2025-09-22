package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/server"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/models"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	redis      *redis.Client
	mu         sync.RWMutex
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	fleetID  string
	userType string // "fleet_manager", "driver", etc.
}

type Message struct {
	Type      string      `json:"type"`
	FleetID   string      `json:"fleet_id,omitempty"`
	VehicleID string      `json:"vehicle_id,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin for now
	},
}

func main() {
	// Initialize base server with common setup
	baseServer, err := server.NewBaseServer("websocket")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize server")
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logrus.WithError(err).Warn("Redis connection failed, continuing without pub/sub")
		redisClient = nil
	}

	// Create WebSocket hub
	hub := &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		redis:      redisClient,
	}

	// Start hub
	go hub.run()

	// Start Redis subscriber if available
	if redisClient != nil {
		go hub.subscribeToRedis(ctx)
	}

	// WebSocket endpoint
	setupWebSocketRoutes(baseServer, hub)

	// Start server
	port := getEnv("WEBSOCKET_PORT", "8083")
	if err := baseServer.Start(port); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}

	// Wait for shutdown
	baseServer.WaitForShutdown()
}

func setupWebSocketRoutes(server *server.BaseServer, hub *Hub) {
	server.Router.GET("/ws", func(c *gin.Context) {
		handleWebSocket(hub, c.Writer, c.Request)
	})

	// Health check endpoint
	server.Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":           "healthy",
			"service":          "websocket",
			"connected_clients": len(hub.clients),
			"timestamp":        time.Now(),
		})
	})

	logrus.Info("WebSocket endpoints configured")
}

func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to upgrade connection")
		return
	}

	// Get client parameters from query string
	fleetID := r.URL.Query().Get("fleet_id")
	userType := r.URL.Query().Get("user_type")

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		fleetID:  fleetID,
		userType: userType,
	}

	client.hub.register <- client

	// Start goroutines for this client
	go client.writePump()
	go client.readPump()
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

			logrus.WithFields(logrus.Fields{
				"fleet_id":  client.fleetID,
				"user_type": client.userType,
				"total_clients": len(h.clients),
			}).Info("Client connected")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

			logrus.WithFields(logrus.Fields{
				"fleet_id":  client.fleetID,
				"user_type": client.userType,
				"total_clients": len(h.clients),
			}).Info("Client disconnected")

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) subscribeToRedis(ctx context.Context) {
	if h.redis == nil {
		return
	}

	// Subscribe to various channels
	pubsub := h.redis.Subscribe(ctx,
		"risk_events",
		"alerts",
		"vehicle_updates",
		"driver_updates",
	)
	defer pubsub.Close()

	logrus.Info("Started Redis subscription for real-time events")

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			logrus.WithError(err).Error("Redis subscription error")
			time.Sleep(time.Second)
			continue
		}

		// Broadcast message to WebSocket clients
		h.broadcast <- []byte(msg.Payload)
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Error("WebSocket error")
			}
			break
		}

		// Handle incoming messages (ping, subscribe to specific updates, etc.)
		logrus.WithField("message", string(message)).Debug("Received WebSocket message")
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// BroadcastRiskEvent sends a risk event to relevant clients
func (h *Hub) BroadcastRiskEvent(event *models.RiskEvent) {
	message := Message{
		Type:      "risk_event",
		FleetID:   "",  // Would get from vehicle relationship
		VehicleID: fmt.Sprintf("%d", event.VehicleID),
		Data:      event,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal risk event")
		return
	}

	h.broadcast <- data
}

// BroadcastAlert sends an alert to relevant clients
func (h *Hub) BroadcastAlert(alert *models.Alert) {
	message := Message{
		Type:      "alert",
		FleetID:   fmt.Sprintf("%d", alert.FleetID),
		Data:      alert,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal alert")
		return
	}

	h.broadcast <- data
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}