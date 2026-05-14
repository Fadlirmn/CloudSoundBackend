package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/sumbul/music-player-backend/internal/delivery/http"
	"github.com/sumbul/music-player-backend/internal/repository"
	"github.com/sumbul/music-player-backend/internal/service"
	"github.com/sumbul/music-player-backend/pkg/external_api"
)

func main() {
	log.Println("Memulai aplikasi...")
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Firebase Firestore
	client, err := repository.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to Firebase: %v", err)
	}
	defer client.Close()

	// Initialize Repositories
	userRepo := repository.NewUserRepository(client)
	playlistRepo := repository.NewPlaylistRepository(client)
	musicRepo := repository.NewMusicRepository(client)

	// Initialize External Clients
	jamendoClient := external_api.NewJamendoClient(os.Getenv("JAMENDO_CLIENT_ID"))

	// Initialize Services
	authService := service.NewAuthService(userRepo)
	musicService := service.NewMusicService(jamendoClient, musicRepo, userRepo)
	playlistService := service.NewPlaylistService(playlistRepo)
	keepAliveService := service.NewKeepAliveService(client)

	// Start Background Workers
	keepAliveService.StartBackgroundWorker()

	// Initialize Handlers
	authHandler := http.NewAuthHandler(authService)
	musicHandler := http.NewMusicHandler(musicService)
	playlistHandler := http.NewPlaylistHandler(playlistService)

	// Setup Router
	r := gin.Default()

	// CORS Middleware
	envOrigins := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string

	if envOrigins != "" {
		allowedOrigins = strings.Split(envOrigins, ",")
	}

	// Selalu tambahkan localhost agar development tetap lancar
	allowedOrigins = append(allowedOrigins, "http://localhost:5173")

	// Hapus spasi jika ada user yang input "url1, url2"
	for i, v := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(v)
	}
	
	allowedMethods := strings.Split(os.Getenv("ALLOWED_METHODS"), ",")
	if len(allowedMethods) == 0 || allowedMethods[0] == "" {
		allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}

	allowedHeaders := strings.Split(os.Getenv("ALLOWED_HEADERS"), ",")
	if len(allowedHeaders) == 0 || allowedHeaders[0] == "" {
		allowedHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}
	
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
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
			meGroup.GET("/profile", authHandler.GetProfile)
			meGroup.GET("/playlists", playlistHandler.GetMyPlaylists)
			meGroup.POST("/playlists", playlistHandler.Create)
			meGroup.POST("/playlists/:id/tracks", playlistHandler.AddTrack)
			
			// Music Activity
			meGroup.POST("/recent", musicHandler.SaveRecentlyPlayed)
			meGroup.POST("/like", musicHandler.ToggleLike)
			meGroup.GET("/liked", musicHandler.GetLikedTracks)
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
