package dto

import "time"

/*************** REQUEST DTOs ***************/

// type UpdateScoreRequest struct {
// 	Delta int `json:"delta" binding:"required"`
// }

// type SetScoreRequest struct {
// 	Score int `json:"score" binding:"required"`
// }

/*************** RESPONSE DTOs ***************/

type LeaderboardEntryResponse struct {
	User          User       `json:"user"`
	Score         int        `json:"score"`
	CurrentStreak int        `json:"current_streak"`
	LongestStreak int        `json:"longest_streak"`
	LastCheckedIn *time.Time `json:"last_checkedin,omitempty"`
}
