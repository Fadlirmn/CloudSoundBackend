package repository

import (
	"fmt"
	"os"

	"github.com/sumbul/music-player-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migration
	err = db.AutoMigrate(
		&models.User{},
		&models.Playlist{},
		&models.Track{},
		&models.RecentlyPlayed{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
