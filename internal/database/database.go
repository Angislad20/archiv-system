package database

import (
	"archiv-system/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=root dbname=archiv_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&models.User{}, &models.Document{})
	if err != nil {
		return nil
	}
	DB = db
	return DB
}
