package service

import (
	"github.com/sumbul/music-player-backend/internal/models"
	"github.com/sumbul/music-player-backend/internal/repository"
)

type PlaylistService interface {
	CreatePlaylist(userID string, title, description string, isPublic bool) (*models.Playlist, error)
	GetUserPlaylists(userID string) ([]models.Playlist, error)
	GetPlaylistByID(id uint) (*models.Playlist, error)
	AddTrackToPlaylist(playlistID uint, track models.Track) error
	RemoveTrackFromPlaylist(playlistID uint, trackTitle string) error
}

type playlistService struct {
	repo repository.PlaylistRepository
}

func NewPlaylistService(repo repository.PlaylistRepository) PlaylistService {
	return &playlistService{repo}
}

func (s *playlistService) CreatePlaylist(userID string, title, description string, isPublic bool) (*models.Playlist, error) {
	playlist := &models.Playlist{
		UserID:      userID,
		Title:       title,
		Description: description,
		IsPublic:    isPublic,
	}
	err := s.repo.Create(playlist)
	return playlist, err
}

func (s *playlistService) GetUserPlaylists(userID string) ([]models.Playlist, error) {
	return s.repo.GetByUserID(userID)
}

func (s *playlistService) GetPlaylistByID(id uint) (*models.Playlist, error) {
	return s.repo.GetByID(id)
}

func (s *playlistService) AddTrackToPlaylist(playlistID uint, track models.Track) error {
	return s.repo.AddTrack(playlistID, &track)
}

func (s *playlistService) RemoveTrackFromPlaylist(playlistID uint, trackTitle string) error {
	return s.repo.RemoveTrack(playlistID, trackTitle)
}
