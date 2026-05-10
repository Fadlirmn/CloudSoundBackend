package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sumbul/music-player-backend/internal/models"
	"github.com/sumbul/music-player-backend/internal/service"
)

type MusicHandler struct {
	service service.MusicService
}

func NewMusicHandler(service service.MusicService) *MusicHandler {
	return &MusicHandler{service}
}

func (h *MusicHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	tracks, err := h.service.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

func (h *MusicHandler) GetTrack(c *gin.Context) {
	id := c.Param("id")
	track, err := h.service.GetTrack(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}

	c.JSON(http.StatusOK, track)
}

func (h *MusicHandler) GetFeed(c *gin.Context) {
	tracks, err := h.service.GetHomeFeed()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

func (h *MusicHandler) GetRecommendations(c *gin.Context) {
	tracks, err := h.service.GetRecommendations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

func (h *MusicHandler) GetMostPlayed(c *gin.Context) {
	tracks, err := h.service.GetMostPlayed()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

func (h *MusicHandler) SaveRecentlyPlayed(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var track models.Track
	if err := c.ShouldBindJSON(&track); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SaveRecentlyPlayed(userID, &track)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *MusicHandler) ToggleLike(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var track models.Track
	if err := c.ShouldBindJSON(&track); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isLiked, err := h.service.ToggleLike(userID, &track)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_liked": isLiked})
}

func (h *MusicHandler) GetLikedTracks(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tracks, err := h.service.GetLikedTracks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}
