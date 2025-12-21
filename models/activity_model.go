package models

import (
	"fmt"
	"time"

	"klubRanks/db"
)

type ActivityLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	ClubID       uint      `gorm:"not null;index" json:"club_id"`
	Action       string    `gorm:"not null" json:"action"`
	UpdatedScore int       `gorm:"not null;default:0" json:"updated_score"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type dailyUserScore struct {
	UserID uint
	Total  int
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}

const (
	ActionJoin   = "joined"
	ActionLeave  = "left"
	ActionUpdate = "update"
)

func AddActivityLog(userID, clubID uint, updatedScore int, action string) error {
	club, err := getClubByID(clubID)
	if err != nil {
		return err
	}
	user, _ := GetUserByID(userID)

	if action == ActionJoin || action == ActionLeave {
		message := Message{
			UserID:    userID,
			ClubID:    clubID,
			Message:   fmt.Sprintf("%s has %s %s.", user.Username, action, club.Name),
			Timestamp: time.Now(),
			Type:      MessageTypeSystem,
		}
		message.AddMessage()

		log := ActivityLog{
			UserID:       userID,
			ClubID:       clubID,
			Action:       action,
			UpdatedScore: 0,
			CreatedAt:    time.Now(),
		}

		return db.DB.Create(&log).Error
	} else {

		message := Message{
			UserID:    userID,
			ClubID:    clubID,
			Message:   user.Username + " increased their count by " + fmt.Sprint(updatedScore) + " " + club.Action,
			Timestamp: time.Now(),
			Type:      MessageTypeSystem,
		}
		message.AddMessage()

		log := ActivityLog{
			UserID:       userID,
			ClubID:       clubID,
			Action:       club.Action,
			UpdatedScore: updatedScore,
			CreatedAt:    time.Now(),
		}

		return db.DB.Create(&log).Error
	}
}

func GetDailyScoresForClub(
	clubID uint,
	day time.Time,
	currentUserID uint,
) (map[string]int, error) {

	startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var rows []dailyUserScore
	err := db.DB.Raw(`
		SELECT user_id, SUM(updated_score) AS total
		FROM activity_logs
		WHERE club_id = ?
		  AND created_at >= ?
		  AND created_at < ?
		GROUP BY user_id
		ORDER BY total DESC
	`, clubID, startOfDay, endOfDay).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	selected := make(map[uint]int)
	for i, r := range rows {
		if i < 3 || r.UserID == currentUserID {
			selected[r.UserID] = r.Total
		}
	}

	result := make(map[string]int)
	result["You"] = 0 // Default value if current user has no activity

	for userID, score := range selected {
		if userID == currentUserID {
			result["You"] = score
			continue
		}

		user, err := GetUserByID(userID)
		if err != nil {
			return nil, err
		}
		result[user.Username] = score
	}

	return result, nil
}
