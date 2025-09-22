package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Features FeatureFlags
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port      string
	Host      string
	Env       string
	JWTSecret string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// FeatureFlags holds feature flag configuration
type FeatureFlags struct {
	EnableRealTimeProcessing bool
	EnableMLRiskScoring      bool
	EnableTelemetrySimulation bool
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:      getEnv("API_PORT", "8080"),
			Host:      getEnv("API_HOST", "0.0.0.0"),
			Env:       getEnv("ENV", "development"),
			JWTSecret: getEnv("JWT_SECRET", ""),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "fleet"),
			Password: getEnv("DB_PASSWORD", "devpass"),
			Database: getEnv("DB_NAME", "fleet_dev"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Features: FeatureFlags{
			EnableRealTimeProcessing:  getEnvAsBool("ENABLE_REAL_TIME_PROCESSING", true),
			EnableMLRiskScoring:       getEnvAsBool("ENABLE_ML_RISK_SCORING", true),
			EnableTelemetrySimulation: getEnvAsBool("ENABLE_TELEMETRY_SIMULATION", true),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}