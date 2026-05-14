package models

import (
	"time"
)

type User struct {
	ID           string    `json:"id" firestore:"id"`
	Name         string    `json:"name" firestore:"name"`
	Email        string    `json:"email" firestore:"email"`
	PasswordHash string    `json:"-" firestore:"password_hash"`
	LastSeen     time.Time `json:"last_seen" firestore:"last_seen"`
	CreatedAt    time.Time `json:"created_at" firestore:"created_at"`
}

type Playlist struct {
	ID          string    `json:"id" firestore:"id"`
	UserID      string    `json:"user_id" firestore:"user_id"`
	Title       string    `json:"title" firestore:"title"`
	Description string    `json:"description" firestore:"description"`
	IsPublic    bool      `json:"is_public" firestore:"is_public"`
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
	Tracks      []Track   `json:"tracks" firestore:"tracks"`
}

type Track struct {
	ID          string `json:"id" firestore:"id"`
	ExternalID  string `json:"external_id" firestore:"external_id"` 
	Title       string `json:"title" firestore:"title"`
	ArtistName  string `json:"artist_name" firestore:"artist_name"`
	AlbumName   string `json:"album_name" firestore:"album_name"`
	Duration    int    `json:"duration" firestore:"duration"`
	AudioURL    string `json:"audio_url" firestore:"audio_url"`
	ImageURL    string `json:"image_url" firestore:"image_url"`
}

type PlaylistTrack struct {
	PlaylistID string    `json:"playlist_id" firestore:"playlist_id"`
	TrackTitle string    `json:"track_title" firestore:"track_title"`
	AddedAt    time.Time `json:"added_at" firestore:"added_at"`
}

type RecentlyPlayed struct {
	UserID     string    `json:"user_id" firestore:"user_id"`
	TrackTitle string    `json:"track_title" firestore:"track_title"`
	PlayedAt   time.Time `json:"played_at" firestore:"played_at"`
}

type LikedTrack struct {
	UserID     string    `json:"user_id" firestore:"user_id"`
	TrackTitle string    `json:"track_title" firestore:"track_title"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
}

type APIUsage struct {
	Datetime string `json:"datetime" firestore:"datetime"`
	Hits     int    `json:"hits" firestore:"hits"`
}

type SystemLog struct {
	ID        string    `json:"id" firestore:"id"`
	Event     string    `json:"event" firestore:"event"`
	Message   string    `json:"message" firestore:"message"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
}
