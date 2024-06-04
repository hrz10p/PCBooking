package models

import (
	"gorm.io/gorm"
	"time"
)

type Computer struct {
	gorm.Model
	Number int    `gorm:"unique;not null"`
	Status string `gorm:"type:varchar(20);default:'available'"` // available, booked
}
type Booking struct {
	gorm.Model
	UserID     string    `gorm:"type:varchar(255);not null"`
	ComputerID uint      `gorm:"not null"`
	StartTime  time.Time `gorm:"not null"`
	EndTime    time.Time `gorm:"not null"`
}
