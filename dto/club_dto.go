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

/*************** RESPONSE DTOs ***************/

type ClubResponse struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	IsPrivate       bool      `json:"is_private"`
	NumberOfMembers int       `json:"number_of_members"`
	CreatedBy       int64     `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}

type MemberResponse struct {
	User     User      `json:"user"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type UserStatsDTO struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	AvatarID string `json:"avatar_id,omitempty"`

	Score int `json:"score"`

	CurrentStreak int `json:"current_streak"`
	LongestStreak int `json:"longest_streak"`

	LastCheckedIn *time.Time `json:"last_checkedin,omitempty"`
	Rank          int        `json:"rank"`
}
