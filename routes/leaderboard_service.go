package routes

import (
	"klubRanks/dto"
	"klubRanks/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddUserToLeaderboard godoc
// @Summary Add user to leaderboard
// @Description Creates leaderboard entry with score 0
// @Tags Leaderboard
// @Security BearerAuth
// @Produce json
// @Param clubId path int true "Club ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/leaderboard/join [post]
func AddUserToLeaderboard(c *gin.Context) {
	clubID, err := strconv.ParseInt(c.Param("clubId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid club id"})
		return
	}

	userID := c.GetInt64("userId")

	if err := models.AddUserToLeaderboard(userID, clubID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "user added to leaderboard",
	})
}

// UpdateLeaderboardScore godoc
// @Summary Update leaderboard score
// @Description Increment user's score
// @Tags Leaderboard
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param clubId path int true "Club ID"
// @Param score body dto.UpdateScoreRequest true "Score delta"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/leaderboard/score [patch]
func UpdateLeaderboardScore(c *gin.Context) {
	clubID, err := strconv.ParseInt(c.Param("clubId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid club id"})
		return
	}

	var req dto.UpdateScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	userID := c.GetInt64("userId")

	err = models.UpdateLeaderboardScore(userID, clubID, req.Delta)
	if err != nil {
		if err.Error() == "leaderboard entry not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "score updated",
	})
}

// SetLeaderboardScore godoc
// @Summary Set leaderboard score
// @Description Set user's score directly (admin)
// @Tags Leaderboard
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param clubId path int true "Club ID"
// @Param score body dto.SetScoreRequest true "Score value"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/leaderboard/score [put]
func SetLeaderboardScore(c *gin.Context) {
	clubID, err := strconv.ParseInt(c.Param("clubId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid club id"})
		return
	}

	var req dto.SetScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	userID := c.GetInt64("userId")

	if err := models.SetLeaderboardScore(userID, clubID, req.Score); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "score set successfully",
	})
}

// GetLeaderboard godoc
// @Summary Get club leaderboard
// @Description Fetch top N users sorted by score
// @Tags Leaderboard
// @Security BearerAuth
// @Produce json
// @Param clubId path int true "Club ID"
// @Param limit query int false "Result limit" default(50)
// @Success 200 {array} dto.LeaderboardEntryResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/leaderboard [get]
func GetLeaderboard(c *gin.Context) {
	clubID, err := strconv.ParseInt(c.Param("clubId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid club id"})
		return
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	entries, err := models.GetLeaderboardForClub(clubID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]dto.LeaderboardEntryResponse, 0, len(entries))
	for _, e := range entries {
		resp = append(resp, dto.LeaderboardEntryResponse{
			UserID:        e.UserID,
			ClubID:        e.ClubID,
			Score:         e.Score,
			LastCheckedIn: e.LastCheckedIn,
		})
	}

	c.JSON(http.StatusOK, resp)
}
