package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	usages, err := h.service.GetRecommendations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usages)
}

func (h *MusicHandler) GetMostPlayed(c *gin.Context) {
	tracks, err := h.service.GetMostPlayed()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}
