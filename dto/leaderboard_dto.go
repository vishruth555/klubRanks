package dto

import "time"

/*************** REQUEST DTOs ***************/

type UpdateScoreRequest struct {
	Delta int `json:"delta" binding:"required"`
}

type SetScoreRequest struct {
	Score int `json:"score" binding:"required"`
}

/*************** RESPONSE DTOs ***************/

type LeaderboardEntryResponse struct {
	UserID        int64      `json:"user_id"`
	ClubID        int64      `json:"club_id"`
	Score         int        `json:"score"`
	LastCheckedIn *time.Time `json:"last_checkedin,omitempty"`
}
