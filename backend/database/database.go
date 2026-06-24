package database

import (
	"log"
	"time"

	"xinhang-backend/config"
	"xinhang-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxConns)
	sqlDB.SetMaxIdleConns(cfg.DBIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	log.Printf("Database connected (maxOpen=%d, maxIdle=%d)", cfg.DBMaxConns, cfg.DBIdleConns)

	if err = DB.AutoMigrate(
		&models.User{},
		&models.Application{},
		&models.News{},
		&models.Notification{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")
}
