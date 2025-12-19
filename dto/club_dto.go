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
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	IsPrivate       bool      `json:"is_private"`
	NumberOfMembers int       `json:"number_of_members"`
	CreatedBy       uint      `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}

type MemberResponse struct {
	User     User      `json:"user"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type GraphDataPoint struct {
	Day    string `json:"name"`
	You    int    `json:"You"`
	Leader int    `json:"Leader"`
}

type UserStats struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	AvatarID string `json:"avatar_id,omitempty"`

	Score int `json:"score"`

	CurrentStreak int `json:"current_streak"`
	LongestStreak int `json:"longest_streak"`

	LastCheckedIn *time.Time `json:"last_checkedin,omitempty"`
	Rank          int        `json:"rank"`

	Percentile string           `json:"percentile"`
	GraphData  []GraphDataPoint `json:"graph_data"`
}
