package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = Migrate(db)
	assert.NoError(t, err)

	return db
}

func TestFleetModel(t *testing.T) {
	db := setupTestDB(t)

	fleet := Fleet{
		Name:         "Test Fleet",
		CompanyName:  "Test Company",
		ContactEmail: "test@example.com",
		Status:       "active",
	}

	// Test Create
	err := db.Create(&fleet).Error
	assert.NoError(t, err)
	assert.NotZero(t, fleet.ID)
	assert.NotZero(t, fleet.CreatedAt)

	// Test Read
	var retrievedFleet Fleet
	err = db.First(&retrievedFleet, fleet.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, fleet.Name, retrievedFleet.Name)
	assert.Equal(t, fleet.CompanyName, retrievedFleet.CompanyName)

	// Test Update
	fleet.Status = "inactive"
	err = db.Save(&fleet).Error
	assert.NoError(t, err)

	err = db.First(&retrievedFleet, fleet.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "inactive", retrievedFleet.Status)
}

func TestDriverModel(t *testing.T) {
	db := setupTestDB(t)

	// Create a fleet first
	fleet := Fleet{
		Name:   "Test Fleet",
		Status: "active",
	}
	err := db.Create(&fleet).Error
	assert.NoError(t, err)

	driver := Driver{
		FleetID:      fleet.ID,
		FirstName:    "John",
		LastName:     "Doe",
		LicenseNum:    "DL123456",
		Email:        "john.doe@example.com",
		Phone:        "555-1234",
		Status:       "active",
		RiskScore:    85.5,
	}

	// Test Create
	err = db.Create(&driver).Error
	assert.NoError(t, err)
	assert.NotZero(t, driver.ID)

	// Test Read with Fleet preload
	var retrievedDriver Driver
	err = db.Preload("Fleet").First(&retrievedDriver, driver.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, driver.FirstName, retrievedDriver.FirstName)
	assert.Equal(t, driver.LastName, retrievedDriver.LastName)
	assert.Equal(t, fleet.Name, retrievedDriver.Fleet.Name)

	// Test Update
	driver.RiskScore = 90.0
	err = db.Save(&driver).Error
	assert.NoError(t, err)

	err = db.First(&retrievedDriver, driver.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, 90.0, retrievedDriver.RiskScore)
}

func TestVehicleModel(t *testing.T) {
	db := setupTestDB(t)

	// Create a fleet and driver first
	fleet := Fleet{
		Name:   "Test Fleet",
		Status: "active",
	}
	err := db.Create(&fleet).Error
	assert.NoError(t, err)

	driver := Driver{
		FleetID:   fleet.ID,
		FirstName: "Jane",
		LastName:  "Smith",
		Status:    "active",
	}
	err = db.Create(&driver).Error
	assert.NoError(t, err)

	vehicle := Vehicle{
		FleetID:        fleet.ID,
		DriverID:       &driver.ID,
		VIN:           "1HGCM82633A123456",
		Make:          "Honda",
		Model:         "Accord",
		Year:          2023,
		LicensePlate:  "ABC123",
		Status:        "active",
	}

	// Test Create
	err = db.Create(&vehicle).Error
	assert.NoError(t, err)
	assert.NotZero(t, vehicle.ID)

	// Test Read with associations
	var retrievedVehicle Vehicle
	err = db.Preload("Fleet").Preload("Driver").First(&retrievedVehicle, vehicle.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, vehicle.VIN, retrievedVehicle.VIN)
	assert.Equal(t, fleet.Name, retrievedVehicle.Fleet.Name)
	assert.Equal(t, driver.FirstName, retrievedVehicle.Driver.FirstName)

	// Test nullable driver
	vehicle2 := Vehicle{
		FleetID:      fleet.ID,
		VIN:         "1HGCM82633A654321",
		Make:        "Toyota",
		Model:       "Camry",
		Year:        2023,
		LicensePlate: "XYZ789",
		Status:      "active",
	}
	err = db.Create(&vehicle2).Error
	assert.NoError(t, err)
	assert.Nil(t, vehicle2.DriverID)
}

func TestTelemetryEventModel(t *testing.T) {
	db := setupTestDB(t)

	// Create necessary relationships
	fleet := Fleet{Name: "Test Fleet", Status: "active"}
	err := db.Create(&fleet).Error
	assert.NoError(t, err)

	vehicle := Vehicle{
		FleetID: fleet.ID,
		VIN:    "1HGCM82633A123456",
		Make:   "Honda",
		Model:  "Accord",
		Year:   2023,
		Status: "active",
	}
	err = db.Create(&vehicle).Error
	assert.NoError(t, err)

	lat := 37.7749
	lon := -122.4194
	speed := 55.0
	accel := 2.5

	event := TelemetryEvent{
		VehicleID:    vehicle.ID,
		EventType:    "location",
		Timestamp:    time.Now(),
		Latitude:     &lat,
		Longitude:    &lon,
		Speed:        &speed,
		Acceleration: &accel,
		Data:         `{"engine_status": "on"}`,
	}

	// Test Create
	err = db.Create(&event).Error
	assert.NoError(t, err)
	assert.NotZero(t, event.ID)

	// Test Read
	var retrievedEvent TelemetryEvent
	err = db.Preload("Vehicle").First(&retrievedEvent, event.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, event.EventType, retrievedEvent.EventType)
	assert.NotNil(t, retrievedEvent.Latitude)
	assert.Equal(t, lat, *retrievedEvent.Latitude)
	assert.Equal(t, vehicle.VIN, retrievedEvent.Vehicle.VIN)

	// Test ProcessedAt field
	assert.Nil(t, retrievedEvent.ProcessedAt)
	now := time.Now()
	retrievedEvent.ProcessedAt = &now
	err = db.Save(&retrievedEvent).Error
	assert.NoError(t, err)
}

func TestRiskEventModel(t *testing.T) {
	db := setupTestDB(t)

	// Create necessary relationships
	fleet := Fleet{Name: "Test Fleet", Status: "active"}
	err := db.Create(&fleet).Error
	assert.NoError(t, err)

	driver := Driver{
		FleetID:   fleet.ID,
		FirstName: "Test",
		LastName:  "Driver",
		Status:    "active",
	}
	err = db.Create(&driver).Error
	assert.NoError(t, err)

	vehicle := Vehicle{
		FleetID:  fleet.ID,
		DriverID: &driver.ID,
		VIN:     "1HGCM82633A123456",
		Make:    "Honda",
		Model:   "Accord",
		Year:    2023,
		Status:  "active",
	}
	err = db.Create(&vehicle).Error
	assert.NoError(t, err)

	lat := 37.7749
	lon := -122.4194

	riskEvent := RiskEvent{
		VehicleID:   vehicle.ID,
		DriverID:    &driver.ID,
		EventType:   "speeding",
		Severity:    "high",
		RiskScore:   85.0,
		Timestamp:   time.Now(),
		Latitude:    &lat,
		Longitude:   &lon,
		Description: "Vehicle exceeded speed limit by 20 mph",
		Data:        `{"speed": 85, "limit": 65}`,
	}

	// Test Create
	err = db.Create(&riskEvent).Error
	assert.NoError(t, err)
	assert.NotZero(t, riskEvent.ID)

	// Test Read with associations
	var retrievedRiskEvent RiskEvent
	err = db.Preload("Vehicle").Preload("Driver").First(&retrievedRiskEvent, riskEvent.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, riskEvent.EventType, retrievedRiskEvent.EventType)
	assert.Equal(t, riskEvent.Severity, retrievedRiskEvent.Severity)
	assert.Equal(t, vehicle.VIN, retrievedRiskEvent.Vehicle.VIN)
	assert.Equal(t, driver.FirstName, retrievedRiskEvent.Driver.FirstName)
}

func TestAlertModel(t *testing.T) {
	db := setupTestDB(t)

	// Create necessary relationships
	fleet := Fleet{Name: "Test Fleet", Status: "active"}
	err := db.Create(&fleet).Error
	assert.NoError(t, err)

	alert := Alert{
		FleetID:  fleet.ID,
		Type:     "risk",
		Priority: "high",
		Title:    "High Risk Event Detected",
		Message:  "Multiple speeding violations detected",
		Status:   "unread",
	}

	// Test Create
	err = db.Create(&alert).Error
	assert.NoError(t, err)
	assert.NotZero(t, alert.ID)

	// Test Read
	var retrievedAlert Alert
	err = db.Preload("Fleet").First(&retrievedAlert, alert.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, alert.Title, retrievedAlert.Title)
	assert.Equal(t, alert.Status, retrievedAlert.Status)
	assert.Equal(t, fleet.Name, retrievedAlert.Fleet.Name)

	// Test Update status
	retrievedAlert.Status = "read"
	err = db.Save(&retrievedAlert).Error
	assert.NoError(t, err)

	var updatedAlert Alert
	err = db.First(&updatedAlert, alert.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "read", updatedAlert.Status)
}

func TestDriverScoreModel(t *testing.T) {
	db := setupTestDB(t)

	// Create necessary relationships
	fleet := Fleet{Name: "Test Fleet", Status: "active"}
	err := db.Create(&fleet).Error
	assert.NoError(t, err)

	driver := Driver{
		FleetID:   fleet.ID,
		FirstName: "Test",
		LastName:  "Driver",
		Status:    "active",
	}
	err = db.Create(&driver).Error
	assert.NoError(t, err)

	score := DriverScore{
		DriverID:        driver.ID,
		OverallScore:    85.5,
		SafetyScore:     90.0,
		EfficiencyScore: 81.0,
		TotalMiles:      1500.5,
		TotalTrips:      120,
		RiskEvents:      5,
		LastUpdated:     time.Now(),
	}

	// Test Create
	err = db.Create(&score).Error
	assert.NoError(t, err)
	assert.NotZero(t, score.ID)

	// Test Read
	var retrievedScore DriverScore
	err = db.Preload("Driver").First(&retrievedScore, score.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, score.OverallScore, retrievedScore.OverallScore)
	assert.Equal(t, score.TotalTrips, retrievedScore.TotalTrips)
	assert.Equal(t, driver.FirstName, retrievedScore.Driver.FirstName)

	// Test Update
	retrievedScore.OverallScore = 88.0
	retrievedScore.RiskEvents = 4
	err = db.Save(&retrievedScore).Error
	assert.NoError(t, err)

	var updatedScore DriverScore
	err = db.First(&updatedScore, score.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, 88.0, updatedScore.OverallScore)
	assert.Equal(t, 4, updatedScore.RiskEvents)
}