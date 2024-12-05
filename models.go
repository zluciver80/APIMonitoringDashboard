package main

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type APIInfo struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"index:,unique"`
	Version     string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type APIMetric struct {
	ID        uint      `gorm:"primaryKey"`
	APIInfoID uint      `gorm:"index"`
	Status    int
	Latency   float64
	Timestamp time.Time
}

func EstablishDatabaseConnection() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&APIInfo{}, &APIMetric{}); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	database, err := EstablishDatabaseConnection()
	if err != nil {
		panic("Failed to connect to the database.")
	}
	database.Create(&APIInfo{Name: "Sample API", Version: "v1", Description: "This is a sample API"})
	database.Create(&APIMetric{APIInfoID: 1, Status: 200, Latency: 123.4, Timestamp: time.Now()})
}