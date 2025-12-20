package dto

import "time"

/*************** REQUEST DTOs ***************/

type CreateClubRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	IsPrivate   bool    `json:"is_private"`
	Action      string  `json:"action" binding:"required"`
}

type UpdateClubRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	IsPrivate   bool    `json:"is_private"`
	Action      string  `json:"action"`
}

/*************** RESPONSE DTOs ***************/

type ClubResponse struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	Code            int       `json:"code"`
	Action          string    `json:"action"`
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
	Day    string         `json:"day"`
	Scores map[string]int `json:"scores"`
}

type UserStats struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	AvatarID string `json:"avatar_id,omitempty"`

	Score int `json:"score"`

	CurrentStreak int `json:"current_streak"`
	LongestStreak int `json:"longest_streak"`

	LastCheckedIn *time.Time `json:"last_checkedin,omitempty"`
	NextCheckIn   *time.Time `json:"next_checkin,omitempty"`
	Rank          int        `json:"rank"`

	GraphData []GraphDataPoint `json:"graph_data"`
}
