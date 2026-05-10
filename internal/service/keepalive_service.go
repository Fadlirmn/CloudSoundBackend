package service

import (
	"fmt"
	"log"
	"time"

	"github.com/sumbul/music-player-backend/internal/models"
	"gorm.io/gorm"
)

type KeepAliveService interface {
	StartBackgroundWorker()
	PerformUpdate()
}

type keepAliveService struct {
	db *gorm.DB
}

func NewKeepAliveService(db *gorm.DB) KeepAliveService {
	return &keepAliveService{db: db}
}

func (s *keepAliveService) StartBackgroundWorker() {
	// Jalankan setiap 24 jam untuk menjaga Supabase tetap aktif
	ticker := time.NewTicker(24 * time.Hour)
	
	log.Println("Background worker untuk Keep-Alive Supabase dimulai (24 jam interval)")
	
	go func() {
		// Jalankan sekali saat startup
		s.PerformUpdate()
		
		for range ticker.C {
			s.PerformUpdate()
		}
	}()
}

func (s *keepAliveService) PerformUpdate() {
	log.Println("Menjalankan tugas harian: Update log dan aktivitas...")

	// 1. Update Log Sistem
	systemLog := models.SystemLog{
		Event:   "DAILY_KEEP_ALIVE",
		Message: fmt.Sprintf("System maintenance & keep-alive executed at %s", time.Now().Format(time.RFC1123)),
	}
	
	if err := s.db.Create(&systemLog).Error; err != nil {
		log.Printf("Error saat membuat system log: %v", err)
	}

	// 2. Hitung Akun Aktif (Hari ini)
	var activeUsers int64
	today := time.Now().Truncate(24 * time.Hour)
	s.db.Model(&models.User{}).Where("last_seen >= ?", today).Count(&activeUsers)

	// 3. Update Ringkasan Aktivitas harian
	var totalRecentlyPlayed int64
	s.db.Model(&models.RecentlyPlayed{}).Where("played_at >= ?", today).Count(&totalRecentlyPlayed)

	activityLog := models.SystemLog{
		Event:   "DAILY_ACTIVITY_SUMMARY",
		Message: fmt.Sprintf("Active Users: %d, Tracks Played today: %d", activeUsers, totalRecentlyPlayed),
	}
	
	if err := s.db.Create(&activityLog).Error; err != nil {
		log.Printf("Error saat membuat activity log: %v", err)
	}

	log.Printf("Keep-alive selesai: %d user aktif, %d trek diputar.", activeUsers, totalRecentlyPlayed)
}
