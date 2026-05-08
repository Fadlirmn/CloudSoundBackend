# Music Player API (Backend)

A high-performance music streaming API built with Go, Gin, and GORM. This backend integrates with the Jamendo API to provide music search, discovery, and personalized library features.

## 🚀 Tech Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin Gonic
- **ORM**: GORM (PostgreSQL)
- **API Integration**: Jamendo Music API
- **Authentication**: JWT (JSON Web Tokens)
- **Environment**: Dotenv for configuration

## ✨ Features

- **Auth**: JWT-based Login & Registration.
- **Music**: 
  - Search tracks from Jamendo.
  - Discovery feed & Recommendations.
  - Most played track tracking.
- **Library**: 
  - Create & manage personal playlists.
  - Add tracks to playlists.
  - Liked/Favorite tracks (managed via frontend/playlists).
- **CORS**: Configurable Allowed Origins for secure deployment.

## 🛠️ Local Setup

1. **Clone the repository**
2. **Setup Environment Variables**:
   Create a `.env` file in the root directory:
   ```env
   PORT=8080
   DB_HOST=localhost
   DB_USER=your_user
   DB_PASSWORD=your_password
   DB_NAME=music_player
   DB_PORT=5432
   JAMENDO_CLIENT_ID=your_jamendo_id
   ALLOWED_ORIGINS=http://localhost:5173
   JWT_SECRET=your_secret_key
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```
4. **Run the server**:
   ```bash
   go run cmd/main.go
   ```

## 🌐 API Endpoints

- `POST /v1/auth/register` - User registration
- `POST /v1/auth/login` - User login
- `GET /v1/music/search?q=...` - Search music
- `GET /v1/music/feed` - Discovery feed
- `GET /v1/music/recommendations` - Activity-based recommendations
- `GET /v1/me/playlists` - Get user playlists
- `POST /v1/me/playlists` - Create new playlist
- `POST /v1/me/playlists/:id/tracks` - Add track to playlist

## 🚢 Deployment

Recommended platforms: **Render**, **Railway**, or **Fly.io**. Ensure you set the `ALLOWED_ORIGINS` environment variable to your production frontend URL.
