package models

import (
	"time"
)

type User struct {
	ID           string    `gorm:"primaryKey;type:varchar(26)" json:"id"`
	Name         string    `json:"name"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	PasswordHash string    `json:"-"`
	LastSeen     time.Time `json:"last_seen"`
	CreatedAt    time.Time `json:"created_at"`
}

type Playlist struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	Tracks      []Track   `gorm:"many2many:playlist_tracks;" json:"tracks"`
}

type Track struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ExternalID  string `gorm:"uniqueIndex" json:"external_id"` // ID from Jamendo/Deezer
	Title       string `json:"title"`
	ArtistName  string `json:"artist_name"`
	AlbumName   string `json:"album_name"`
	Duration    int    `json:"duration"`
	AudioURL    string `json:"audio_url"`
	ImageURL    string `json:"image_url"`
}

type PlaylistTrack struct {
	PlaylistID      uint      `gorm:"primaryKey"`
	TrackExternalID string    `gorm:"primaryKey"`
	AddedAt         time.Time `json:"added_at"`
}

type RecentlyPlayed struct {
	UserID          string    `gorm:"primaryKey" json:"user_id"`
	TrackExternalID string    `gorm:"primaryKey" json:"track_external_id"`
	PlayedAt        time.Time `json:"played_at"`
}

type LikedTrack struct {
	UserID          string    `gorm:"primaryKey" json:"user_id"`
	TrackExternalID string    `gorm:"primaryKey" json:"track_external_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type APIUsage struct {
	Datetime string `json:"datetime"`
	Hits     int    `json:"hits"`
}

type SystemLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Event     string    `json:"event"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
