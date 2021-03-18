package config

import (
	"app/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	database, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	database.AutoMigrate(&models.Product{})
	// database.Migrator().DropTable(&models.Products{})

	db = database
}

func GetDB() *gorm.DB {
	return db
}
