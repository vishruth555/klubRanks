package dto

import "time"

/*************** REQUEST DTOs ***************/

type SendMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

/*************** RESPONSE DTOs ***************/

type ClubMessageResponse struct {
	MessageID int64     `json:"message_id"`
	UserID    int64     `json:"user_id"`
	ClubID    int64     `json:"club_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
