package models

import (
	"time"

	"gorm.io/gorm"
)

// Vehicle represents a fleet vehicle
type Vehicle struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	VIN         string    `json:"vin" gorm:"uniqueIndex;size:17"`
	Make        string    `json:"make"`
	Model       string    `json:"model"`
	Year        int       `json:"year"`
	LicensePlate string   `json:"license_plate"`
	FleetID     uint      `json:"fleet_id"`
	Fleet       Fleet     `json:"fleet"`
	DriverID    *uint     `json:"driver_id"`
	Driver      *Driver   `json:"driver,omitempty"`
	Status      string    `json:"status" gorm:"default:active"` // active, maintenance, inactive
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Driver represents a vehicle driver
type Driver struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	EmployeeID  string    `json:"employee_id" gorm:"uniqueIndex"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email" gorm:"uniqueIndex"`
	Phone       string    `json:"phone"`
	LicenseNum  string    `json:"license_number"`
	FleetID     uint      `json:"fleet_id"`
	Fleet       Fleet     `json:"fleet"`
	Status      string    `json:"status" gorm:"default:active"` // active, suspended, inactive
	RiskScore   float64   `json:"risk_score" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Fleet represents a fleet organization
type Fleet struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	CompanyName string    `json:"company_name"`
	ContactEmail string   `json:"contact_email"`
	Status      string    `json:"status" gorm:"default:active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TelemetryEvent represents raw telemetry data from vehicles
type TelemetryEvent struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	VehicleID   uint      `json:"vehicle_id"`
	Vehicle     Vehicle   `json:"vehicle"`
	EventType   string    `json:"event_type"` // location, speed, acceleration, harsh_braking, etc.
	Timestamp   time.Time `json:"timestamp"`
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`
	Speed       *float64  `json:"speed"`       // mph
	Acceleration *float64 `json:"acceleration"` // m/s²
	Data        string    `json:"data" gorm:"type:json"` // Additional event-specific data
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// RiskEvent represents detected risky behavior
type RiskEvent struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	VehicleID   uint      `json:"vehicle_id"`
	Vehicle     Vehicle   `json:"vehicle"`
	DriverID    *uint     `json:"driver_id"`
	Driver      *Driver   `json:"driver,omitempty"`
	EventType   string    `json:"event_type"` // speeding, harsh_braking, rapid_acceleration, fatigue
	Severity    string    `json:"severity"`   // low, medium, high, critical
	RiskScore   float64   `json:"risk_score"` // 0-100
	Timestamp   time.Time `json:"timestamp"`
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`
	Description string    `json:"description"`
	Data        string    `json:"data" gorm:"type:json"`
	Status      string    `json:"status" gorm:"default:open"` // open, acknowledged, resolved
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Alert represents system-generated alerts
type Alert struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FleetID     uint      `json:"fleet_id"`
	Fleet       Fleet     `json:"fleet"`
	VehicleID   *uint     `json:"vehicle_id"`
	Vehicle     *Vehicle  `json:"vehicle,omitempty"`
	DriverID    *uint     `json:"driver_id"`
	Driver      *Driver   `json:"driver,omitempty"`
	RiskEventID *uint     `json:"risk_event_id"`
	RiskEvent   *RiskEvent `json:"risk_event,omitempty"`
	Type        string    `json:"type"`     // risk, maintenance, system
	Priority    string    `json:"priority"` // low, medium, high, critical
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	Status      string    `json:"status" gorm:"default:unread"` // unread, read, dismissed
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DriverScore represents aggregated driver performance metrics
type DriverScore struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	DriverID       uint      `json:"driver_id" gorm:"uniqueIndex"`
	Driver         Driver    `json:"driver"`
	OverallScore   float64   `json:"overall_score"`   // 0-100
	SafetyScore    float64   `json:"safety_score"`    // 0-100
	EfficiencyScore float64  `json:"efficiency_score"` // 0-100
	TotalMiles     float64   `json:"total_miles"`
	TotalTrips     int       `json:"total_trips"`
	RiskEvents     int       `json:"risk_events"`
	LastUpdated    time.Time `json:"last_updated"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Migrate runs the database migrations
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Fleet{},
		&Driver{},
		&Vehicle{},
		&TelemetryEvent{},
		&RiskEvent{},
		&Alert{},
		&DriverScore{},
	)
}