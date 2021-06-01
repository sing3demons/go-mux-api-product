package config

import (
	"fmt"
	"github/sing3demons/go_mux_api/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Bangkok", host, user, pass, name)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	database.AutoMigrate(&models.Product{})
	database.AutoMigrate(&models.User{})
	// database.Migrator().DropTable(&models.Product{})

	db = database
}

func GetDB() *gorm.DB {
	return db
}
