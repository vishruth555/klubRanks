package routes

import (
	"klubRanks/dto"
	"klubRanks/logger"
	"klubRanks/models"
	"klubRanks/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Signup godoc
// @Summary Create a new user
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.SignupRequest true "User signup payload"
// @Success 200 {object} dto.MessageResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /signup [post]
func signup(context *gin.Context) {
	var req dto.SignupRequest
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "could not parse data"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password,
		AvatarID: req.AvatarID,
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Could not create user " + err.Error()})
		return
	}
	context.JSON(http.StatusOK, dto.MessageResponse{Message: "user created successfully"})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.LoginRequest true "User login payload"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /login [post]
func login(context *gin.Context) {
	var req dto.LoginRequest
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "could not parse data"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password,
		AvatarID: req.AvatarID,
	}

	logger.LogDebug("Attempting login for user: " + user.Username)

	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	logger.LogDebug("User " + user.Username + " authenticated successfully with ID " + strconv.FormatUint(uint64(user.ID), 10) + " and AvatarID " + user.AvatarID)

	token, err := utils.GenerateToken(user.Username, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Return the token AND the user details
	context.JSON(http.StatusOK, dto.LoginResponse{
		Message: "login successful",
		Token:   token,
		User: dto.User{
			ID:       user.ID,
			Username: user.Username,
			AvatarID: user.AvatarID,
		},
	})
}

// UpdateAvatar godoc
// @Summary Update user avatar
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param avatar body dto.UpdateAvatarRequest true "Avatar payload"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/avatar [put]
func UpdateAvatar(c *gin.Context) {
	var req dto.UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
		return
	}

	userID := c.GetUint("userId")

	if err := models.UpdateAvatar(userID, req.AvatarID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "avatar updated"})
}
