package models

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestCreateComputer(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Computer{})
	assert.NoError(t, err)

	computer := Computer{Number: 1, Status: "available"}
	result := db.Create(&computer)
	assert.NoError(t, result.Error)
	assert.Equal(t, uint(1), computer.ID)
}

func TestCreateBooking(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Computer{}, &Booking{})
	assert.NoError(t, err)

	computer := Computer{Number: 1, Status: "available"}
	db.Create(&computer)

	stTime, err := time.Parse(time.DateTime, "2024-06-01 10:00:00")
	assert.NoError(t, err)

	endTime, err := time.Parse(time.DateTime, "2024-06-01 12:00:00")
	assert.NoError(t, err)

	booking := Booking{
		UserID:     "user123",
		ComputerID: computer.ID,
		StartTime:  stTime,
		EndTime:    endTime,
	}
	result := db.Create(&booking)
	assert.NoError(t, result.Error)
	assert.Equal(t, uint(1), booking.ID)
}
