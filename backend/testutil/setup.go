package testutil

import (
	"xinhang-backend/database"
	"xinhang-backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to open test database: " + err.Error())
	}
	db.AutoMigrate(&models.User{}, &models.Application{})
	database.DB = db
}

func CleanDB() {
	database.DB.Exec("DELETE FROM applications")
	database.DB.Exec("DELETE FROM users")
}
