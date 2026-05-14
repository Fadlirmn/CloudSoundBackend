package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/sumbul/music-player-backend/internal/models"
	"google.golang.org/api/iterator"
)

type PlaylistRepository interface {
	Create(playlist *models.Playlist) error
	GetByUserID(userID string) ([]models.Playlist, error)
	GetByID(id string) (*models.Playlist, error)
	AddTrack(playlistID string, track *models.Track) error
	RemoveTrack(playlistID string, trackTitle string) error
}

type playlistRepo struct {
	client *firestore.Client
}

func NewPlaylistRepository(client *firestore.Client) PlaylistRepository {
	return &playlistRepo{client}
}

func (r *playlistRepo) Create(playlist *models.Playlist) error {
	ctx := context.Background()
	docRef, _, err := r.client.Collection("playlists").Add(ctx, playlist)
	if err != nil {
		return err
	}
	playlist.ID = docRef.ID
	// Update the ID in the document as well
	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "id", Value: docRef.ID},
	})
	return err
}

func (r *playlistRepo) GetByUserID(userID string) ([]models.Playlist, error) {
	ctx := context.Background()
	var playlists []models.Playlist
	iter := r.client.Collection("playlists").Where("user_id", "==", userID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var playlist models.Playlist
		if err := doc.DataTo(&playlist); err != nil {
			return nil, err
		}
		playlist.ID = doc.Ref.ID
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (r *playlistRepo) GetByID(id string) (*models.Playlist, error) {
	ctx := context.Background()
	doc, err := r.client.Collection("playlists").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var playlist models.Playlist
	if err := doc.DataTo(&playlist); err != nil {
		return nil, err
	}
	playlist.ID = doc.Ref.ID
	return &playlist, nil
}

func (r *playlistRepo) AddTrack(playlistID string, track *models.Track) error {
	ctx := context.Background()
	docRef := r.client.Collection("playlists").Doc(playlistID)
	
	// In Firestore, we can just append to the Tracks array
	_, err := docRef.Update(ctx, []firestore.Update{
		{
			Path:  "tracks",
			Value: firestore.ArrayUnion(track),
		},
	})
	return err
}

func (r *playlistRepo) RemoveTrack(playlistID string, trackTitle string) error {
	ctx := context.Background()
	docRef := r.client.Collection("playlists").Doc(playlistID)
	
	// Removing from array in Firestore by value is tricky if we only have the title.
	// We'll need to get the playlist, find the track, and then remove it.
	doc, err := docRef.Get(ctx)
	if err != nil {
		return err
	}
	
	var playlist models.Playlist
	if err := doc.DataTo(&playlist); err != nil {
		return err
	}
	
	var trackToRemove *models.Track
	for _, t := range playlist.Tracks {
		if t.Title == trackTitle {
			trackToRemove = &t
			break
		}
	}
	
	if trackToRemove == nil {
		return nil // Not found
	}
	
	_, err = docRef.Update(ctx, []firestore.Update{
		{
			Path:  "tracks",
			Value: firestore.ArrayRemove(trackToRemove),
		},
	})
	return err
}
