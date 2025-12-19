package routes

import (
	"klubRanks/dto"
	"klubRanks/models"
	"klubRanks/utils"
	"net/http"

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

	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

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
