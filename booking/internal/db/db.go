package db

import (
	"booking/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func InitDB(cfg config.Config) (*gorm.DB, error) {
	var err error
	DB, err := gorm.Open(postgres.Open(cfg.PostgresURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return DB, err
}
