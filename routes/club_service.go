package routes

import (
	"klubRanks/dto"
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

	c.JSON(http.StatusCreated, dto.ClubResponse{
		ID:          club.ID,
		Name:        club.Name,
		Description: club.Description,
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
		resp = append(resp, dto.ClubResponse{
			ID:          club.ID,
			Name:        club.Name,
			Description: club.Description,
			IsPrivate:   club.IsPrivate,
			CreatedBy:   club.CreatedBy,
			CreatedAt:   club.CreatedAt,
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
// @Param member body dto.AddMemberRequest true "Member info"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/members [post]
func AddMember(c *gin.Context) {
	clubID, err := strconv.Atoi(c.Param("clubId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid club id"})
		return
	}

	var req dto.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	if err := models.AddMember(int(req.UserID), clubID, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

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
		resp = append(resp, dto.MemberResponse{
			ID:       m.ID,
			UserID:   m.UserID,
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}
