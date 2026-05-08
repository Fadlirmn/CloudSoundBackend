package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/sumbul/music-player-backend/internal/delivery/http"
	"github.com/sumbul/music-player-backend/internal/repository"
	"github.com/sumbul/music-player-backend/internal/service"
	"github.com/sumbul/music-player-backend/pkg/external_api"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Database
	db, err := repository.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)

	// Initialize External Clients
	jamendoClient := external_api.NewJamendoClient(os.Getenv("JAMENDO_CLIENT_ID"))

	// Initialize Services
	authService := service.NewAuthService(userRepo)
	musicService := service.NewMusicService(jamendoClient)
	playlistService := service.NewPlaylistService(playlistRepo)

	// Initialize Handlers
	authHandler := http.NewAuthHandler(authService)
	musicHandler := http.NewMusicHandler(musicService)
	playlistHandler := http.NewPlaylistHandler(playlistService)

	// Setup Router
	r := gin.Default()

	// CORS Middleware
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173"
	}
	
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	v1 := r.Group("/v1")
	{
		// Auth
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Music Proxy
		musicGroup := v1.Group("/music")
		{
			musicGroup.GET("/search", musicHandler.Search)
			musicGroup.GET("/tracks/:id", musicHandler.GetTrack)
			musicGroup.GET("/feed", musicHandler.GetFeed)
			musicGroup.GET("/recommendations", musicHandler.GetRecommendations)
			musicGroup.GET("/most-played", musicHandler.GetMostPlayed)
		}

		// User Library (Protected)
		meGroup := v1.Group("/me")
		meGroup.Use(http.AuthMiddleware())
		{
			meGroup.GET("/playlists", playlistHandler.GetMyPlaylists)
			meGroup.POST("/playlists", playlistHandler.Create)
			meGroup.POST("/playlists/:id/tracks", playlistHandler.AddTrack)
		}
	}

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
