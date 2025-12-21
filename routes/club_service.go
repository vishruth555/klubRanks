package routes

import (
	"klubRanks/config"
	"klubRanks/dto"
	"klubRanks/logger"
	"klubRanks/models"
	"net/http"
	"strconv"
	"time"

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

	userID := c.GetUint("userId")

	club := models.Club{
		Name:        req.Name,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
		Action:      req.Action,
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
	}

	c.JSON(http.StatusCreated, dto.ClubResponse{
		ID:              club.ID,
		Name:            club.Name,
		Description:     club.Description,
		Code:            club.Code,
		IsPrivate:       club.IsPrivate,
		Action:          club.Action,
		NumberOfMembers: 1,
		CreatedBy:       club.CreatedBy,
		CreatedAt:       club.CreatedAt,
	})
}

// UpdateClub godoc
// @Summary Update club details
// @Tags Clubs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param clubId path int true "Club ID"
// @Param club body dto.UpdateClubRequest true "Update club payload"
// @Success 200 {object} dto.ClubResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId} [put]
func UpdateClub(c *gin.Context) {
	clubID, err := strconv.ParseUint(c.Param("clubId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid club id"})
		return
	}

	var req dto.UpdateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	userID := c.GetUint("userId")

	// Verify ownership
	club, err := models.GetClub(uint(clubID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "club not found"})
		return
	}

	if club.CreatedBy != userID {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "only the creator can edit this club"})
		return
	}

	// Update fields
	club.Name = req.Name
	club.Description = req.Description
	club.IsPrivate = req.IsPrivate
	club.Action = req.Action

	if err := club.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ClubResponse{
		ID:          club.ID,
		Name:        club.Name,
		Description: club.Description,
		Code:        club.Code,
		Action:      club.Action,
		IsPrivate:   club.IsPrivate,
		CreatedBy:   club.CreatedBy,
		CreatedAt:   club.CreatedAt,
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
	userID := c.GetUint("userId")

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
		rank, err := models.GetUserRankInClub(userID, club.ID)
		stats, err := models.GetLeaderboardEntryForUser(userID, club.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		var nextCheckIn *time.Time
		if stats.LastCheckedIn != nil {
			t := stats.LastCheckedIn.Add(time.Duration(config.AppConfig.Server.CoolDownMinutes) * time.Minute)
			nextCheckIn = &t
		}

		resp = append(resp, dto.ClubResponse{
			ID:              club.ID,
			Name:            club.Name,
			Description:     club.Description,
			Code:            club.Code,
			Action:          club.Action,
			IsPrivate:       club.IsPrivate,
			NumberOfMembers: int(numberOfMembers),
			LastCheckedIn:   stats.LastCheckedIn,
			NextCheckIn:     nextCheckIn,
			CurrentRank:     rank,
			CreatedBy:       club.CreatedBy,
			CreatedAt:       club.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// JoinClub godoc (Renamed from AddMember)
// @Summary Join a club using invite code
// @Tags Clubs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param code path string true "Club Invite Code"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/join/{code} [post]
func JoinClub(c *gin.Context) {
	clubCode := c.Param("code") // Changed from clubId to code

	userID := c.GetUint("userId")
	role := "member"

	if err := models.AddMember(userID, clubCode, role); err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	logger.LogInfo("User", userID, "joined club with code", clubCode)

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "joined club successfully",
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
	clubID, err := strconv.ParseUint(c.Param("clubId"), 10, 64)

	members, err := models.GetClubMembers(uint(clubID))
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

// @Summary Leave a club
// @Tags Clubs
// @Security BearerAuth
// @Produce json
// @Param clubId path int true "Club ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/members [delete]
func LeaveClub(c *gin.Context) {
	clubID, err := strconv.ParseInt(c.Param("clubId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid club id"})
		return
	}

	userID := c.GetUint("userId")

	if err := models.RemoveMember(userID, uint(clubID)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "left club successfully",
	})
}

// GetClubMembers godoc
// @Summary Get club user stats for current user
// @Tags Clubs
// @Security BearerAuth
// @Produce json
// @Param clubId path int true "Club ID"
// @Success 200 {array} dto.UserStats
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/stats/me [get]
func GetCurrentUserStats(c *gin.Context) {
	clubID, _ := strconv.ParseInt(c.Param("clubId"), 10, 64)

	userID := c.GetUint("userId")

	userStats, err := getClubUserStats(userID, uint(clubID))
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
	clubID, _ := strconv.ParseUint(c.Param("clubId"), 10, 64)

	userID, _ := strconv.ParseUint(c.Param("userId"), 10, 64)

	userStats, err := getClubUserStats(uint(userID), uint(clubID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, userStats)
}

func getClubUserStats(userID uint, clubID uint) (dto.UserStats, error) {

	var userStats dto.UserStats

	user, err := models.GetUserByID(userID)
	if err != nil {
		return userStats, err
	}

	logger.LogInfo("Fetching stats for user:", userID, "in club:", clubID)

	stats, err := models.GetLeaderboardEntryForUser(userID, clubID)
	if err != nil {
		return userStats, err
	}

	rank, err := models.GetUserRankInClub(userID, clubID)
	if err != nil {
		return userStats, err
	}

	// 🔥 Build graph data using daily scores
	graphData := make([]dto.GraphDataPoint, 0)
	now := time.Now()

	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dayLabel := date.Format("Mon")

		scores, err := models.GetDailyScoresForClub(clubID, date, userID)
		if err != nil {
			return userStats, err
		}

		graphData = append(graphData, dto.GraphDataPoint{
			Day:    dayLabel,
			Scores: scores,
		})
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
		GraphData:     graphData,
	}

	return userStats, nil
}
