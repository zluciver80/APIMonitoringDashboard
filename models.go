package main

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// API describes the basic information of an API
type API struct {
	ID          uint      `gorm:"primaryKey"` // ID is the primary key
	Name        string    `gorm:"index:,unique"` // Name is a unique identifier for the API
	Version     string    // Version of the API
	Description string    // Description provides more details about the API
	CreatedAt   time.Time // CreatedAt records when the API entry was created
	UpdatedAt   time.Time // UpdatedAt records the last update time
}

// APIMetric records performance metrics of an API
type APIMetric struct {
	ID        uint      `gorm:"primaryKey"` // ID is the primary key
	APIID     uint      `gorm:"index"` // APIID links to the corresponding API entry
	Status    int       // Status code returned by the API call
	Latency   float64   // Latency in milliseconds of the API call
	Timestamp time.Time // Timestamp when the metric was recorded
}

// ConnectToDatabase initializes the database connection and migrates the schema
func ConnectToDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate the API and APIMetric models
	if err := database.AutoMigrate(&API{}, &APIMetric{}); err != nil {
		return nil, err
	}

	return database, nil
}

func main() {
	db, err := ConnectToDatabase()
	if err != nil {
		panic("Failed to connect to the database")
	}
	// Creating a sample API entry
	db.Create(&API{Name: "Sample API", Version: "v1", Description: "This is a sample API"})

	// Creating a sample API metric entry
	db.Create(&APIMetric{APIID: 1, Status: 200, Latency: 123.4, Timestamp: time.Now()})
}