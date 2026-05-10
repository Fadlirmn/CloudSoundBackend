package repository

import (
	"time"
	"github.com/sumbul/music-player-backend/internal/models"
	"gorm.io/gorm"
)

type MusicRepository interface {
	SaveRecentlyPlayed(userID string, track *models.Track) error
	GetRecentlyPlayed(userID string, limit int) ([]models.RecentlyPlayed, error)
	ToggleLike(userID string, track *models.Track) (bool, error)
	GetLikedTracks(userID string) ([]models.Track, error)
}

type musicRepo struct {
	db *gorm.DB
}

func NewMusicRepository(db *gorm.DB) MusicRepository {
	return &musicRepo{db}
}

func (r *musicRepo) ensureTrackExists(track *models.Track) error {
	var existingTrack models.Track
	err := r.db.Where("external_id = ?", track.ExternalID).First(&existingTrack).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.db.Create(track).Error
		}
		return err
	}
	return nil
}

func (r *musicRepo) SaveRecentlyPlayed(userID string, track *models.Track) error {
	if err := r.ensureTrackExists(track); err != nil {
		return err
	}

	recent := models.RecentlyPlayed{
		UserID:          userID,
		TrackExternalID: track.ExternalID,
		PlayedAt:        time.Now(),
	}
	// Upsert
	return r.db.Save(&recent).Error
}

func (r *musicRepo) GetRecentlyPlayed(userID string, limit int) ([]models.RecentlyPlayed, error) {
	var results []models.RecentlyPlayed
	err := r.db.Where("user_id = ?", userID).Order("played_at desc").Limit(limit).Find(&results).Error
	return results, err
}

func (r *musicRepo) ToggleLike(userID string, track *models.Track) (bool, error) {
	if err := r.ensureTrackExists(track); err != nil {
		return false, err
	}

	var liked models.LikedTrack
	err := r.db.Where("user_id = ? AND track_external_id = ?", userID, track.ExternalID).First(&liked).Error
	
	if err == nil {
		err = r.db.Delete(&liked).Error
		return false, err
	}
	
	if err == gorm.ErrRecordNotFound {
		liked = models.LikedTrack{
			UserID:          userID,
			TrackExternalID: track.ExternalID,
			CreatedAt:       time.Now(),
		}
		err = r.db.Create(&liked).Error
		return true, err
	}
	
	return false, err
}

func (r *musicRepo) GetLikedTracks(userID string) ([]models.Track, error) {
	var tracks []models.Track
	err := r.db.Table("tracks").
		Joins("join liked_tracks on liked_tracks.track_external_id = tracks.external_id").
		Where("liked_tracks.user_id = ?", userID).
		Order("liked_tracks.created_at desc").
		Find(&tracks).Error
	return tracks, err
}
