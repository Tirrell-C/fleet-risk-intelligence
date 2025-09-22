package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/config"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/database"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/models"
)

type RiskEngine struct {
	db     *gorm.DB
	config *config.Config
}

type RiskAnalyzer struct {
	SpeedThreshold       float64
	AccelerationThreshold float64
	BrakingThreshold     float64
}

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

	engine := &RiskEngine{
		db:     db,
		config: cfg,
	}

	analyzer := &RiskAnalyzer{
		SpeedThreshold:        getEnvAsFloat("SPEED_THRESHOLD", 80.0),       // mph
		AccelerationThreshold: getEnvAsFloat("ACCEL_THRESHOLD", 4.0),        // m/s²
		BrakingThreshold:     getEnvAsFloat("BRAKING_THRESHOLD", -6.0),      // m/s²
	}

	// Start background risk processing
	go engine.startRiskProcessing(analyzer)

	// Start driver score calculation
	go engine.startDriverScoreCalculation()

	logrus.Info("Risk engine started - processing telemetry data")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Risk engine shutting down...")
}

// startRiskProcessing continuously processes telemetry data for risk detection
func (re *RiskEngine) startRiskProcessing(analyzer *RiskAnalyzer) {
	ticker := time.NewTicker(30 * time.Second) // Process every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			re.processUnprocessedTelemetry(analyzer)
		}
	}
}

// processUnprocessedTelemetry finds and analyzes new telemetry events
func (re *RiskEngine) processUnprocessedTelemetry(analyzer *RiskAnalyzer) {
	var events []models.TelemetryEvent

	// Get unprocessed telemetry events from the last hour
	result := re.db.Where("processed_at IS NULL AND created_at > ?",
		time.Now().Add(-1*time.Hour)).
		Order("timestamp ASC").
		Limit(1000).
		Find(&events)

	if result.Error != nil {
		logrus.WithError(result.Error).Error("Failed to fetch unprocessed telemetry")
		return
	}

	logrus.WithField("count", len(events)).Debug("Processing telemetry events")

	for _, event := range events {
		risks := analyzer.AnalyzeEvent(&event)

		for _, risk := range risks {
			if err := re.createRiskEvent(risk); err != nil {
				logrus.WithError(err).Error("Failed to create risk event")
				continue
			}

			// Create alert if risk is high severity
			if risk.Severity == "high" || risk.Severity == "critical" {
				if err := re.createAlert(risk); err != nil {
					logrus.WithError(err).Error("Failed to create alert")
				}
			}
		}

		// Mark as processed
		now := time.Now()
		re.db.Model(&event).Update("processed_at", &now)
	}
}

// startDriverScoreCalculation periodically updates driver risk scores
func (re *RiskEngine) startDriverScoreCalculation() {
	ticker := time.NewTicker(10 * time.Minute) // Update every 10 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			re.updateDriverScores()
		}
	}
}

// updateDriverScores calculates and updates driver risk scores
func (re *RiskEngine) updateDriverScores() {
	var drivers []models.Driver
	if err := re.db.Where("status = ?", "active").Find(&drivers).Error; err != nil {
		logrus.WithError(err).Error("Failed to fetch active drivers")
		return
	}

	for _, driver := range drivers {
		score := re.calculateDriverScore(driver.ID)

		// Update driver's risk score
		re.db.Model(&driver).Update("risk_score", score.OverallScore)

		// Upsert driver score record
		var existingScore models.DriverScore
		result := re.db.Where("driver_id = ?", driver.ID).First(&existingScore)

		if result.Error != nil {
			// Create new score record
			score.DriverID = driver.ID
			if err := re.db.Create(&score).Error; err != nil {
				logrus.WithError(err).Error("Failed to create driver score")
			}
		} else {
			// Update existing score
			if err := re.db.Model(&existingScore).Updates(&score).Error; err != nil {
				logrus.WithError(err).Error("Failed to update driver score")
			}
		}
	}

	logrus.WithField("drivers", len(drivers)).Info("Updated driver scores")
}

// AnalyzeEvent analyzes a telemetry event for potential risks
func (ra *RiskAnalyzer) AnalyzeEvent(event *models.TelemetryEvent) []models.RiskEvent {
	var risks []models.RiskEvent

	// Speed analysis
	if event.Speed != nil && *event.Speed > ra.SpeedThreshold {
		severity := "medium"
		riskScore := 50.0

		if *event.Speed > ra.SpeedThreshold*1.3 {
			severity = "high"
			riskScore = 75.0
		}
		if *event.Speed > ra.SpeedThreshold*1.5 {
			severity = "critical"
			riskScore = 90.0
		}

		risks = append(risks, models.RiskEvent{
			VehicleID:   event.VehicleID,
			EventType:   "speeding",
			Severity:    severity,
			RiskScore:   riskScore,
			Timestamp:   event.Timestamp,
			Latitude:    event.Latitude,
			Longitude:   event.Longitude,
			Description: fmt.Sprintf("Vehicle exceeded speed limit: %.1f mph", *event.Speed),
			Data:        fmt.Sprintf(`{"speed": %.1f, "threshold": %.1f}`, *event.Speed, ra.SpeedThreshold),
		})
	}

	// Harsh acceleration analysis
	if event.Acceleration != nil && *event.Acceleration > ra.AccelerationThreshold {
		risks = append(risks, models.RiskEvent{
			VehicleID:   event.VehicleID,
			EventType:   "rapid_acceleration",
			Severity:    "medium",
			RiskScore:   60.0,
			Timestamp:   event.Timestamp,
			Latitude:    event.Latitude,
			Longitude:   event.Longitude,
			Description: fmt.Sprintf("Harsh acceleration detected: %.1f m/s²", *event.Acceleration),
			Data:        fmt.Sprintf(`{"acceleration": %.1f, "threshold": %.1f}`, *event.Acceleration, ra.AccelerationThreshold),
		})
	}

	// Harsh braking analysis
	if event.Acceleration != nil && *event.Acceleration < ra.BrakingThreshold {
		risks = append(risks, models.RiskEvent{
			VehicleID:   event.VehicleID,
			EventType:   "harsh_braking",
			Severity:    "medium",
			RiskScore:   65.0,
			Timestamp:   event.Timestamp,
			Latitude:    event.Latitude,
			Longitude:   event.Longitude,
			Description: fmt.Sprintf("Harsh braking detected: %.1f m/s²", *event.Acceleration),
			Data:        fmt.Sprintf(`{"acceleration": %.1f, "threshold": %.1f}`, *event.Acceleration, ra.BrakingThreshold),
		})
	}

	return risks
}

// calculateDriverScore computes comprehensive driver safety metrics
func (re *RiskEngine) calculateDriverScore(driverID uint) models.DriverScore {
	var score models.DriverScore

	// Get risk events from last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	var riskCount int64
	re.db.Model(&models.RiskEvent{}).
		Where("driver_id = ? AND created_at > ?", driverID, thirtyDaysAgo).
		Count(&riskCount)

	// Get total driving metrics (simplified calculation)
	var totalMiles float64 = 1000.0 // Mock data - would calculate from telemetry
	var totalTrips int = 50         // Mock data

	// Calculate scores (0-100 scale)
	safetyScore := math.Max(0, 100.0-(float64(riskCount)*5.0))
	efficiencyScore := 85.0 // Mock efficiency score
	overallScore := (safetyScore + efficiencyScore) / 2.0

	score.OverallScore = overallScore
	score.SafetyScore = safetyScore
	score.EfficiencyScore = efficiencyScore
	score.TotalMiles = totalMiles
	score.TotalTrips = totalTrips
	score.RiskEvents = int(riskCount)
	score.LastUpdated = time.Now()

	return score
}

// createRiskEvent saves a new risk event to the database
func (re *RiskEngine) createRiskEvent(risk models.RiskEvent) error {
	return re.db.Create(&risk).Error
}

// createAlert creates an alert for high-priority risk events
func (re *RiskEngine) createAlert(risk models.RiskEvent) error {
	var vehicle models.Vehicle
	if err := re.db.Preload("Fleet").First(&vehicle, risk.VehicleID).Error; err != nil {
		return err
	}

	alert := models.Alert{
		FleetID:     vehicle.FleetID,
		VehicleID:   &risk.VehicleID,
		DriverID:    risk.DriverID,
		RiskEventID: &risk.ID,
		Type:        "risk",
		Priority:    mapSeverityToPriority(risk.Severity),
		Title:       fmt.Sprintf("%s Alert", formatEventType(risk.EventType)),
		Message:     risk.Description,
		Status:      "unread",
	}

	return re.db.Create(&alert).Error
}

func mapSeverityToPriority(severity string) string {
	switch severity {
	case "critical":
		return "critical"
	case "high":
		return "high"
	case "medium":
		return "medium"
	default:
		return "low"
	}
}

func formatEventType(eventType string) string {
	switch eventType {
	case "speeding":
		return "Speeding"
	case "harsh_braking":
		return "Harsh Braking"
	case "rapid_acceleration":
		return "Rapid Acceleration"
	default:
		return "Risk Event"
	}
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

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}