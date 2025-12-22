package routes

import (
	"klubRanks/dto"
	"klubRanks/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SendMessage godoc
// @Summary Send message to club chat
// @Description Send a chat message inside a club
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param clubId path int true "Club ID"
// @Param message body dto.SendMessageRequest true "Message payload"
// @Success 201 {object} dto.ClubMessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/messages [post]
func SendMessage(c *gin.Context) {
	clubID, err := strconv.ParseUint(c.Param("clubId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid club id",
		})
		return
	}

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	userID := c.GetUint("userId")

	msg := models.Message{
		ClubID:    uint(clubID),
		UserID:    userID,
		Message:   req.Message,
		Type:      models.MessageTypeUser,
		ReplyToID: req.ReplyToID,
	}

	if err := msg.AddMessage(); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, dto.MessageResponse{
		Message: "message sent successfully",
	})
}

// GetClubMessages godoc
// @Summary Get club messages
// @Description Fetch paginated messages for a club
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param clubId path int true "Club ID"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} dto.ClubMessageResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clubs/{clubId}/messages [get]
func GetClubMessages(c *gin.Context) {
	clubID, err := strconv.ParseUint(c.Param("clubId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid club id",
		})
		return
	}

	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	messages, err := models.GetMessagesForClub(uint(clubID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp := make([]dto.ClubMessageResponse, 0, len(messages))
	// for _, m := range messages {
	// 	user, err := models.GetUserByID(m.UserID)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
	// 			Error: err.Error(),
	// 		})
	// 		return
	// 	}
	// 	resp = append(resp, dto.ClubMessageResponse{
	// 		User: dto.User{
	// 			ID:       user.ID,
	// 			Username: user.Username,
	// 			AvatarID: user.AvatarID,
	// 		},
	// 		Type:      m.Type,
	// 		Message:   m.Message,
	// 		Timestamp: m.Timestamp,
	// 	})
	// }
	for _, m := range messages {
		user, err := models.GetUserByID(m.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		userDto := dto.User{
			ID:       user.ID, // Uses preloaded User
			Username: user.Username,
			AvatarID: user.AvatarID,
		}

		var replyToDto *dto.ReplyInfo
		if m.ReplyTo != nil {
			replyToUser, err := models.GetUserByID(m.ReplyTo.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: err.Error(),
				})
				return
			}
			replyToDto = &dto.ReplyInfo{
				User: dto.User{
					ID:       replyToUser.ID,
					Username: replyToUser.Username,
					AvatarID: replyToUser.AvatarID,
				},
				Message: m.ReplyTo.Message,
			}
		}

		resp = append(resp, dto.ClubMessageResponse{
			ID:        m.ID,
			User:      userDto,
			Type:      m.Type,
			Message:   m.Message,
			Timestamp: m.Timestamp,
			ReplyTo:   replyToDto,
		})
	}

	c.JSON(http.StatusOK, resp)
}
