package dto

import "time"

/*************** REQUEST DTOs ***************/

type SendMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

/*************** RESPONSE DTOs ***************/

type ClubMessageResponse struct {
	User      User      `json:"user"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}
