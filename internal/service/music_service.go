package service

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/sumbul/music-player-backend/internal/models"
	"github.com/sumbul/music-player-backend/internal/repository"
	"github.com/sumbul/music-player-backend/pkg/external_api"
)

type MusicService interface {
	Search(query string) ([]models.Track, error)
	GetTrack(id string) (*models.Track, error)
	GetHomeFeed() ([]models.Track, error)
	GetRecommendations() ([]models.Track, error)
	GetMostPlayed() ([]models.Track, error)
	SaveRecentlyPlayed(userID string, track *models.Track) error
	ToggleLike(userID string, track *models.Track) (bool, error)
	GetLikedTracks(userID string) ([]models.Track, error)
}

type musicService struct {
	client   *external_api.JamendoClient
	repo     repository.MusicRepository
	userRepo repository.UserRepository
}

func NewMusicService(client *external_api.JamendoClient, repo repository.MusicRepository, userRepo repository.UserRepository) MusicService {
	return &musicService{client, repo, userRepo}
}

func (s *musicService) Search(query string) ([]models.Track, error) {
	jamendoTracks, err := s.client.SearchTracks(query, 20)
	if err != nil {
		return nil, err
	}

	var tracks []models.Track
	for _, t := range jamendoTracks {
		tracks = append(tracks, models.Track{
			ExternalID: t.ID,
			Title:      t.Name,
			ArtistName: t.ArtistName,
			AlbumName:  t.AlbumName,
			Duration:   t.Duration,
			AudioURL:   t.Audio,
			ImageURL:   t.Image,
		})
	}

	return tracks, nil
}

func (s *musicService) GetTrack(id string) (*models.Track, error) {
	t, err := s.client.GetTrackByID(id)
	if err != nil {
		return nil, err
	}

	return &models.Track{
		ExternalID: t.ID,
		Title:      t.Name,
		ArtistName: t.ArtistName,
		AlbumName:  t.AlbumName,
		Duration:   t.Duration,
		AudioURL:   t.Audio,
		ImageURL:   t.Image,
	}, nil
}

func (s *musicService) GetHomeFeed() ([]models.Track, error) {
	jamendoTracks, err := s.client.GetFeed(20)
	if err != nil {
		return nil, err
	}

	return s.mapJamendoToInternal(jamendoTracks), nil
}

func (s *musicService) GetRecommendations() ([]models.Track, error) {
	jamendoTracks, err := s.client.GetPopularTracks(20)
	if err != nil {
		return nil, err
	}

	return s.mapJamendoToInternal(jamendoTracks), nil
}

func (s *musicService) GetMostPlayed() ([]models.Track, error) {
	jamendoTracks, err := s.client.GetMostPlayedTracks(20)
	if err != nil {
		return nil, err
	}

	return s.mapJamendoToInternal(jamendoTracks), nil
}

func (s *musicService) mapJamendoToInternal(jamendoTracks []external_api.JamendoTrack) []models.Track {
	var tracks []models.Track
	for _, t := range jamendoTracks {
		tracks = append(tracks, models.Track{
			ExternalID: t.ID,
			Title:      t.Name,
			ArtistName: t.ArtistName,
			AlbumName:  t.AlbumName,
			Duration:   t.Duration,
			AudioURL:   t.Audio,
			ImageURL:   t.Image,
		})
	}
	return tracks
}

func (s *musicService) SaveRecentlyPlayed(userID string, track *models.Track) error {
	_ = s.userRepo.UpdateLastSeen(userID)
	return s.repo.SaveRecentlyPlayed(userID, track)
}

func (s *musicService) ToggleLike(userID string, track *models.Track) (bool, error) {
	_ = s.userRepo.UpdateLastSeen(userID)
	return s.repo.ToggleLike(userID, track)
}

func (s *musicService) GetLikedTracks(userID string) ([]models.Track, error) {
	return s.repo.GetLikedTracks(userID)
}
