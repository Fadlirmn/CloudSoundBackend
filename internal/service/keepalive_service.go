package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/sumbul/music-player-backend/internal/models"
	"google.golang.org/api/iterator"
)

type KeepAliveService interface {
	StartBackgroundWorker()
	PerformUpdate()
}

type keepAliveService struct {
	client *firestore.Client
}

func NewKeepAliveService(client *firestore.Client) KeepAliveService {
	return &keepAliveService{client: client}
}

func (s *keepAliveService) StartBackgroundWorker() {
	// Jalankan setiap 24 jam
	ticker := time.NewTicker(24 * time.Hour)
	
	log.Println("Background worker untuk Keep-Alive/Maintenance dimulai (24 jam interval)")
	
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
	ctx := context.Background()

	// 1. Update Log Sistem
	systemLog := models.SystemLog{
		Event:     "DAILY_MAINTENANCE",
		Message:   fmt.Sprintf("System maintenance executed at %s", time.Now().Format(time.RFC1123)),
		CreatedAt: time.Now(),
	}
	
	_, _, err := s.client.Collection("system_logs").Add(ctx, systemLog)
	if err != nil {
		log.Printf("Error saat membuat system log: %v", err)
	}

	// 2. Hitung Akun Aktif (Hari ini)
	today := time.Now().Truncate(24 * time.Hour)
	activeUsersIter := s.client.Collection("users").Where("last_seen", ">=", today).Documents(ctx)
	activeUsers := 0
	for {
		_, err := activeUsersIter.Next()
		if err == iterator.Done {
			break
		}
		if err == nil {
			activeUsers++
		}
	}

	// 3. Update Ringkasan Aktivitas harian
	tracksPlayedIter := s.client.Collection("recently_played").Where("played_at", ">=", today).Documents(ctx)
	totalRecentlyPlayed := 0
	for {
		_, err := tracksPlayedIter.Next()
		if err == iterator.Done {
			break
		}
		if err == nil {
			totalRecentlyPlayed++
		}
	}

	activityLog := models.SystemLog{
		Event:     "DAILY_ACTIVITY_SUMMARY",
		Message:   fmt.Sprintf("Active Users: %d, Tracks Played today: %d", activeUsers, totalRecentlyPlayed),
		CreatedAt: time.Now(),
	}
	
	_, _, err = s.client.Collection("system_logs").Add(ctx, activityLog)
	if err != nil {
		log.Printf("Error saat membuat activity log: %v", err)
	}

	log.Printf("Maintenance selesai: %d user aktif, %d trek diputar.", activeUsers, totalRecentlyPlayed)
}
