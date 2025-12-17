package dto

import "time"

/*************** REQUEST DTOs ***************/

type CreateClubRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	IsPrivate   bool    `json:"is_private"`
}

type UpdateClubRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	IsPrivate   bool    `json:"is_private"`
}

type AddMemberRequest struct {
	UserID int64  `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=admin member"`
}

/*************** RESPONSE DTOs ***************/

type ClubResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	IsPrivate   bool      `json:"is_private"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type MemberResponse struct {
	ID       int64     `json:"id"`
	UserID   int64     `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}
