package db

import (
	"booking/internal/models"
	"gorm.io/gorm"
	"log"
)

func SeedComputers(db *gorm.DB) {
	computers := []models.Computer{
		{Number: 1, Status: "available"},
		{Number: 2, Status: "available"},
		{Number: 3, Status: "available"},
		{Number: 4, Status: "available"},
		{Number: 5, Status: "available"},
		{Number: 6, Status: "available"},
	}

	for _, computer := range computers {
		var existing models.Computer
		if err := db.First(&existing, "number = ?", computer.Number).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&computer).Error; err != nil {
					log.Fatalf("Failed to seed computer: %s", err)
				}
			} else {
				log.Fatalf("Failed to check existing computer: %s", err)
			}
		}
	}
}
