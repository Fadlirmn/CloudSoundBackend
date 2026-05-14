package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sumbul/music-player-backend/internal/models"
	"github.com/sumbul/music-player-backend/internal/service"
)

type PlaylistHandler struct {
	service service.PlaylistService
}

func NewPlaylistHandler(service service.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{service}
}

type createPlaylistRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

func (h *PlaylistHandler) Create(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var req createPlaylistRequest
	if err := h.ShouldBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playlist, err := h.service.CreatePlaylist(userID, req.Title, req.Description, req.IsPublic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, playlist)
}

func (h *PlaylistHandler) GetMyPlaylists(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	playlists, err := h.service.GetUserPlaylists(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, playlists)
}

func (h *PlaylistHandler) AddTrack(c *gin.Context) {
	id := c.Param("id")

	var track models.Track
	if err := h.ShouldBindJSON(c, &track); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddTrackToPlaylist(id, track)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Track added to playlist"})
}

// Helper method to bind JSON
func (h *PlaylistHandler) ShouldBindJSON(c *gin.Context, obj interface{}) error {
	return c.ShouldBindJSON(obj)
}
