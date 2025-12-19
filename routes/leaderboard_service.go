package routes

import (
	"klubRanks/config"
	"klubRanks/dto"
	"klubRanks/logger"
	"klubRanks/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateLeaderboardScore godoc
// @Summary Update leaderboard score
// @Description Increment user's score
// @Tags Leaderboard
// @Security BearerAuth
// @Produce json
// @Param clubId path int true "Club ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/leaderboard/score [post]
func UpdateLeaderboardScore(c *gin.Context) {
	clubID, err := strconv.ParseUint(c.Param("clubId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid club id"})
		return
	}

	userID := c.GetUint("userId")
	logger.LogInfo("Updating leaderboard score for user: ", userID, " in club: ", clubID)

	err = models.UpdateLeaderboardScore(userID, uint(clubID), config.AppConfig.Server.Counter)
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
	clubID, err := strconv.ParseUint(c.Param("clubId"), 10, 64)
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

	entries, err := models.GetLeaderboardForClub(uint(clubID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]dto.LeaderboardEntryResponse, 0, len(entries))
	for _, e := range entries {
		user, err := models.GetUserByID(e.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		resp = append(resp, dto.LeaderboardEntryResponse{
			User: dto.User{
				ID:       user.ID,
				Username: user.Username,
				AvatarID: user.AvatarID,
			},
			CurrentStreak: e.CurrentStreak,
			Score:         e.Score,
			LastCheckedIn: e.LastCheckedIn,
		})
	}

	c.JSON(http.StatusOK, resp)
}
