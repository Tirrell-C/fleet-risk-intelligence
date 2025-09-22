package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/auth"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/models"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/server"
)

type AuthService struct {
	db         *gorm.DB
	jwtManager *auth.JWTManager
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	User      models.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role"`
}

func main() {
	// Initialize base server with common setup
	baseServer, err := server.NewBaseServer("auth")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize server")
	}

	// Initialize JWT manager
	jwtSecret := baseServer.Config.Server.JWTSecret
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
		logrus.Warn("Using default JWT secret - change this in production!")
	}

	jwtManager := auth.NewJWTManager(jwtSecret, 24*time.Hour)

	// Create auth service
	authService := &AuthService{
		db:         baseServer.DB,
		jwtManager: jwtManager,
	}

	// Setup routes
	setupAuthRoutes(baseServer, authService)

	// Start server
	if err := baseServer.Start(baseServer.Config.Server.Port); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}

	// Wait for shutdown
	baseServer.WaitForShutdown()
}

func setupAuthRoutes(server *server.BaseServer, authService *AuthService) {
	api := server.Router.Group("/api/v1/auth")

	// Public routes
	api.POST("/login", authService.login)
	api.POST("/register", authService.register)
	api.POST("/refresh", authService.refreshToken)

	// Protected routes
	authMiddleware := auth.NewAuthMiddleware(authService.jwtManager)
	protected := api.Group("")
	protected.Use(authMiddleware.RequireAuth())
	protected.GET("/me", authService.getProfile)
	protected.PUT("/me", authService.updateProfile)
	protected.POST("/logout", authService.logout)

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(authMiddleware.RequireAuth())
	admin.Use(authMiddleware.RequireRole("super_admin", "fleet_admin"))
	admin.GET("/users", authService.listUsers)
	admin.POST("/users", authService.createUser)
	admin.PUT("/users/:id", authService.updateUser)
	admin.DELETE("/users/:id", authService.deleteUser)

	logrus.Info("Auth endpoints configured")
}

func (s *AuthService) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := s.db.Where("email = ? AND status = ?", req.Email, "active").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Parse fleet IDs
	var fleetIDs []string
	if user.FleetIDs != "" {
		json.Unmarshal([]byte(user.FleetIDs), &fleetIDs)
	}

	// Generate JWT token
	token, err := s.jwtManager.Generate(
		strconv.Itoa(int(user.ID)),
		user.Email,
		user.Role,
		fleetIDs,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	// Create session record
	session := models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.db.Create(&session)

	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		User:      user,
		ExpiresAt: session.ExpiresAt,
	})
}

func (s *AuthService) register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "fleet_manager"
	}

	// Create new user
	user := models.User{
		Email:     req.Email,
		Password:  req.Password, // Will be hashed in BeforeCreate hook
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		Status:    "active",
	}

	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func (s *AuthService) refreshToken(c *gin.Context) {
	// Implementation for refresh token
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (s *AuthService) getProfile(c *gin.Context) {
	claims, exists := auth.GetUserFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *AuthService) updateProfile(c *gin.Context) {
	claims, exists := auth.GetUserFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email" binding:"omitempty,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *AuthService) logout(c *gin.Context) {
	// Delete session token
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		token := authHeader[7:] // Remove "Bearer " prefix
		s.db.Where("token = ?", token).Delete(&models.Session{})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (s *AuthService) listUsers(c *gin.Context) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (s *AuthService) createUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	user := models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		Status:    "active",
	}

	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

func (s *AuthService) updateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email" binding:"omitempty,email"`
		Role      string `json:"role"`
		Status    string `json:"status"`
		FleetIDs  string `json:"fleet_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.FleetIDs != "" {
		user.FleetIDs = req.FleetIDs
	}

	if err := s.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (s *AuthService) deleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}