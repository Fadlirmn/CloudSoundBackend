package repository

import (
	"fmt"
	"os"
	"log"

	"github.com/sumbul/music-player-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	var dsn string
	
	// Cek baris DATABASE_URL dulu
	if url := os.Getenv("DATABASE_URL"); url != "" {
		dsn = url
	} else {
		// Jika tidak ada, baru susun manual
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_SSL_MODE"),
			os.Getenv("DB_TIMEZONE"),
		)
	}

	log.Printf("Menghubungkan ke database...")

	// Tambahkan Config khusus untuk debugging
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // WAJIB untuk port 6543 Supabase
	}), &gorm.Config{
		PrepareStmt: false, // Nonaktifkan prepared statements secara global
	})
	
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi: %v", err)
	}

	// TES KONEKSI (Ini yang bikin hang kalau koneksi lambat)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	log.Println("Melakukan ping ke database...")
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database tidak merespon: %v", err)
	}

	log.Println("Sinkronisasi tabel (AutoMigrate)...")
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Playlist{},
		&models.Track{},
		&models.RecentlyPlayed{},
		&models.LikedTrack{},
		&models.SystemLog{},
	}

	for _, model := range modelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			log.Printf("Warning: Gagal migrasi model %T: %v", model, err)
			// Jangan langsung fatal, coba lanjut ke yang lain
		}
	}
	
	return db, nil
}