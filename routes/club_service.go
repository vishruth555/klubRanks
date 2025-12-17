package routes

import (
	"klubRanks/dto"
	"klubRanks/logger"
	"klubRanks/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateClub godoc
// @Summary Create a new club
// @Description Create a club and add creator as admin
// @Tags Clubs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param club body dto.CreateClubRequest true "Create club payload"
// @Success 201 {object} dto.ClubResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs [post]
func CreateClub(c *gin.Context) {
	var req dto.CreateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	userID := c.GetInt64("userId")

	club := models.Club{
		Name:        req.Name,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
		CreatedBy:   userID,
	}

	if err := club.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	logger.LogInfo("Club created with ID:", club.ID)
	if err := models.AddUserToLeaderboard(club.CreatedBy, club.ID); err != nil {
		logger.LogError("Failed to add user to leaderboard:", err)
		return
	}

	c.JSON(http.StatusCreated, dto.ClubResponse{
		ID:              club.ID,
		Name:            club.Name,
		Description:     club.Description,
		IsPrivate:       club.IsPrivate,
		NumberOfMembers: 1,
		CreatedBy:       club.CreatedBy,
		CreatedAt:       club.CreatedAt,
	})
}

// GetMyClubs godoc
// @Summary Get user's clubs
// @Description Get all clubs the user is a member of
// @Tags Clubs
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.ClubResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs [get]
func GetMyClubs(c *gin.Context) {
	userID := c.GetInt64("userId")

	clubs, err := models.GetClubsForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp := make([]dto.ClubResponse, 0, len(clubs))
	for _, club := range clubs {
		numberOfMembers, err := models.GetMemberCountForClub(club.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		resp = append(resp, dto.ClubResponse{
			ID:              club.ID,
			Name:            club.Name,
			Description:     club.Description,
			IsPrivate:       club.IsPrivate,
			NumberOfMembers: int(numberOfMembers),
			CreatedBy:       club.CreatedBy,
			CreatedAt:       club.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// AddMember godoc
// @Summary Add member to club
// @Tags Clubs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param clubId path int true "Club ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/members [post]
func AddMember(c *gin.Context) {
	clubID, err := strconv.ParseInt(c.Param("clubId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid club id"})
		return
	}

	userID := c.GetInt64("userId")
	role := "member"

	if err := models.AddMember(userID, clubID, role); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	logger.LogInfo("User", userID, "added to club", clubID)

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "member added successfully",
	})
}

// GetClubMembers godoc
// @Summary Get club members
// @Tags Clubs
// @Security BearerAuth
// @Produce json
// @Param clubId path int true "Club ID"
// @Success 200 {array} dto.MemberResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/members [get]
func GetClubMembers(c *gin.Context) {
	clubID, _ := strconv.Atoi(c.Param("clubId"))

	members, err := models.GetClubMembers(clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]dto.MemberResponse, 0, len(members))
	for _, m := range members {
		user, err := models.GetUserByID(m.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
			return
		}
		resp = append(resp, dto.MemberResponse{
			User: dto.User{
				ID:       user.ID,
				Username: user.Username,
				AvatarID: user.AvatarID,
			},
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// GetClubMembers godoc
// @Summary Get club user stats for current user
// @Tags Clubs
// @Security BearerAuth
// @Produce json
// @Param clubId path int true "Club ID"
// @Success 200 {array} dto.UserStatsDTO
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/stats/me [get]
func GetCurrentUserStats(c *gin.Context) {
	clubID, _ := strconv.ParseInt(c.Param("clubId"), 10, 64)

	userID := c.GetInt64("userId")

	userStats, err := getClubUserStats(userID, clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, userStats)
}

// GetClubMembers godoc
// @Summary Get club user stats with id
// @Tags Clubs
// @Security BearerAuth
// @Produce json
// @Param clubId path int true "Club ID"
// @Param userId path int true "User ID"
// @Success 200 {array} dto.UserStats
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/stats/{userId} [get]
func GetUserStats(c *gin.Context) {
	clubID, _ := strconv.ParseInt(c.Param("clubId"), 10, 64)

	userID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)

	userStats, err := getClubUserStats(userID, clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, userStats)
}

func getClubUserStats(userID int64, clubID int64) (dto.UserStats, error) {

	var userStats dto.UserStats
	user, err := models.GetUserByID(userID)
	if err != nil {
		return userStats, err
	}

	logger.LogInfo("Fetching stats for user: ", userID, "in club: ", clubID)

	stats, err := models.GetLeaderboardEntryForUser(userID, clubID)
	if err != nil {
		return userStats, err
	}
	logger.LogDebug("Leaderboard stats: ", stats)
	rank, err := models.GetUserRankInClub(userID, clubID)
	if err != nil {
		return userStats, err
	}

	userStats = dto.UserStats{
		UserID:        user.ID,
		Username:      user.Username,
		AvatarID:      user.AvatarID,
		Score:         stats.Score,
		CurrentStreak: stats.CurrentStreak,
		LongestStreak: stats.LongestStreak,
		LastCheckedIn: stats.LastCheckedIn,
		Rank:          rank,
	}
	return userStats, nil
}
