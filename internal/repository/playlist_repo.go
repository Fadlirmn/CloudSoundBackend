package repository

import (
	"github.com/sumbul/music-player-backend/internal/models"
	"gorm.io/gorm"
)

type PlaylistRepository interface {
	Create(playlist *models.Playlist) error
	GetByUserID(userID string) ([]models.Playlist, error)
	GetByID(id uint) (*models.Playlist, error)
	AddTrack(playlistID uint, track *models.Track) error
	RemoveTrack(playlistID uint, trackExternalID string) error
}

type playlistRepo struct {
	db *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) PlaylistRepository {
	return &playlistRepo{db}
}

func (r *playlistRepo) Create(playlist *models.Playlist) error {
	return r.db.Create(playlist).Error
}

func (r *playlistRepo) GetByUserID(userID string) ([]models.Playlist, error) {
	var playlists []models.Playlist
	err := r.db.Where("user_id = ?", userID).Find(&playlists).Error
	return playlists, err
}

func (r *playlistRepo) GetByID(id uint) (*models.Playlist, error) {
	var playlist models.Playlist
	err := r.db.Preload("Tracks").First(&playlist, id).Error
	if err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (r *playlistRepo) AddTrack(playlistID uint, track *models.Track) error {
	// Ensure track exists in database (caching metadata)
	var existingTrack models.Track
	err := r.db.Where("external_id = ?", track.ExternalID).First(&existingTrack).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := r.db.Create(track).Error; err != nil {
				return err
			}
			existingTrack = *track
		} else {
			return err
		}
	}

	return r.db.Model(&models.Playlist{ID: playlistID}).Association("Tracks").Append(&existingTrack)
}

func (r *playlistRepo) RemoveTrack(playlistID uint, trackExternalID string) error {
	var track models.Track
	err := r.db.Where("external_id = ?", trackExternalID).First(&track).Error
	if err != nil {
		return err
	}
	return r.db.Model(&models.Playlist{ID: playlistID}).Association("Tracks").Delete(&track)
}
