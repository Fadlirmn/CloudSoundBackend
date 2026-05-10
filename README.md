# CloudSound API (Backend)

A high-performance music streaming API built with Go, Gin, and GORM. Optimized for reliability with title-based track indexing and session persistence.

## 🚀 Tech Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin Gonic
- **ORM**: GORM (PostgreSQL)
- **Database**: Supabase (Postgres with Connection Pooling)
- **API Integration**: Jamendo Music API
- **Authentication**: JWT (JSON Web Tokens)

## ✨ Key Features

- **Auth System**:
  - JWT-based Login & Registration.
  - Persistent Session Verification via `/v1/me/profile`.
- **Music Logic**: 
  - **Title-indexed Tracks**: Every track is uniquely identified by its title in the database, preventing sync issues across different API providers.
  - Automated Metadata Caching for user activities.
- **Library Management**: 
  - Playlist creation and track associations.
  - Persistent "Liked Tracks" and "Recently Played" history.
- **System Stability**:
  - Daily Keep-Alive worker to prevent Supabase project pausing.
  - Robust CORS configuration with dynamic environment support.

## 🛠️ Local Setup

1. **Clone the repository**
2. **Setup Environment Variables**:
   Create a `.env` file in the root directory:
   ```env
   PORT=8080
   DATABASE_URL=your_postgres_url
   JWT_SECRET=your_secret_key
   JAMENDO_CLIENT_ID=your_jamendo_id
   ALLOWED_ORIGINS=http://localhost:5173
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

### Auth
- `POST /v1/auth/register` - Register new account
- `POST /v1/auth/login` - Login and get JWT
- `GET /v1/me/profile` - Verify session and get user info (Protected)

### Music
- `GET /v1/music/search?q=...` - Search music
- `GET /v1/music/feed` - Discovery feed
- `GET /v1/music/recommendations` - Popular track recommendations
- `GET /v1/music/most-played` - Global top tracks

### User Library (Protected)
- `GET /v1/me/playlists` - List user playlists
- `POST /v1/me/playlists` - Create playlist
- `POST /v1/me/playlists/:id/tracks` - Add track to playlist (uses Title)
- `POST /v1/me/like` - Toggle like for a track
- `GET /v1/me/liked` - List liked tracks
- `POST /v1/me/recent` - Save play history

## 🚢 Deployment

Optimized for **Render**. Ensure you set the `ALLOWED_ORIGINS` to match your frontend domain to allow cross-origin requests.
