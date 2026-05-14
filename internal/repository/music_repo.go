package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/sumbul/music-player-backend/internal/models"
	"google.golang.org/api/iterator"
)

type MusicRepository interface {
	SaveRecentlyPlayed(userID string, track *models.Track) error
	GetRecentlyPlayed(userID string, limit int) ([]models.RecentlyPlayed, error)
	ToggleLike(userID string, track *models.Track) (bool, error)
	GetLikedTracks(userID string) ([]models.Track, error)
}

type musicRepo struct {
	client *firestore.Client
}

func NewMusicRepository(client *firestore.Client) MusicRepository {
	return &musicRepo{client}
}

func (r *musicRepo) ensureTrackExists(track *models.Track) error {
	ctx := context.Background()
	// Use Title as document ID for tracks to ensure uniqueness and easy lookup
	docRef := r.client.Collection("tracks").Doc(track.Title)
	_, err := docRef.Get(ctx)
	if err != nil {
		// If not exists, create it
		_, err = docRef.Set(ctx, track)
		return err
	}
	return nil
}

func (r *musicRepo) SaveRecentlyPlayed(userID string, track *models.Track) error {
	if err := r.ensureTrackExists(track); err != nil {
		return err
	}

	ctx := context.Background()
	// Use a composite ID for recently played to allow easy updates
	id := userID + "_" + track.Title
	recent := models.RecentlyPlayed{
		UserID:     userID,
		TrackTitle: track.Title,
		PlayedAt:   time.Now(),
	}
	_, err := r.client.Collection("recently_played").Doc(id).Set(ctx, recent)
	return err
}

func (r *musicRepo) GetRecentlyPlayed(userID string, limit int) ([]models.RecentlyPlayed, error) {
	ctx := context.Background()
	var results []models.RecentlyPlayed
	iter := r.client.Collection("recently_played").
		Where("user_id", "==", userID).
		OrderBy("played_at", firestore.Desc).
		Limit(limit).
		Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var recent models.RecentlyPlayed
		if err := doc.DataTo(&recent); err != nil {
			return nil, err
		}
		results = append(results, recent)
	}
	return results, nil
}

func (r *musicRepo) ToggleLike(userID string, track *models.Track) (bool, error) {
	if err := r.ensureTrackExists(track); err != nil {
		return false, err
	}

	ctx := context.Background()
	id := userID + "_" + track.Title
	docRef := r.client.Collection("liked_tracks").Doc(id)
	_, err := docRef.Get(ctx)
	
	if err == nil {
		// Already liked, so unlike
		_, err = docRef.Delete(ctx)
		return false, err
	}
	
	// Not liked, so like
	liked := models.LikedTrack{
		UserID:     userID,
		TrackTitle: track.Title,
		CreatedAt:  time.Now(),
	}
	_, err = docRef.Set(ctx, liked)
	return true, err
}

func (r *musicRepo) GetLikedTracks(userID string) ([]models.Track, error) {
	ctx := context.Background()
	var tracks []models.Track
	
	// First get all liked tracks for this user
	iter := r.client.Collection("liked_tracks").
		Where("user_id", "==", userID).
		OrderBy("created_at", firestore.Desc).
		Documents(ctx)
	
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		
		var liked models.LikedTrack
		if err := doc.DataTo(&liked); err != nil {
			return nil, err
		}
		
		// Then fetch the actual track data
		trackDoc, err := r.client.Collection("tracks").Doc(liked.TrackTitle).Get(ctx)
		if err == nil {
			var track models.Track
			if err := trackDoc.DataTo(&track); err == nil {
				tracks = append(tracks, track)
			}
		}
	}
	return tracks, nil
}
