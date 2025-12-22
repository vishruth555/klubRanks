package dto

import "time"

/*************** REQUEST DTOs ***************/

type SendMessageRequest struct {
	Message   string `json:"message" binding:"required"`
	ReplyToID *uint  `json:"reply_to_id"`
}

/*************** RESPONSE DTOs ***************/

type ReplyInfo struct {
	User    User   `json:"user"`
	Message string `json:"message"`
}

type ClubMessageResponse struct {
	ID        uint       `json:"id"`
	User      User       `json:"user"`
	Message   string     `json:"message"`
	Timestamp time.Time  `json:"timestamp"`
	Type      string     `json:"type"`
	ReplyTo   *ReplyInfo `json:"reply_to,omitempty"`
}
