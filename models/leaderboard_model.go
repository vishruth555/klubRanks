package models

import (
	"fmt"
	"time"

	"klubRanks/db"

	"gorm.io/gorm"
)

type LeaderboardEntry struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `gorm:"not null;index:idx_user_club,unique" json:"user_id"`
	ClubID        uint       `gorm:"not null;index:idx_user_club,unique;index" json:"club_id"`
	Score         int        `gorm:"not null;default:0" json:"score"`
	CurrentStreak int        `gorm:"not null;default:0" json:"current_streak"`
	LongestStreak int        `gorm:"not null;default:0" json:"longest_streak"`
	LastCheckedIn *time.Time `gorm:"column:last_checkedin" json:"last_checkedin,omitempty"`
}

func (LeaderboardEntry) TableName() string {
	return "leaderboard"
}

func AddUserToLeaderboard(userID, clubID uint) error {
	entry := LeaderboardEntry{
		UserID: userID,
		ClubID: clubID,
		Score:  0,
	}
	return db.DB.Create(&entry).Error
}

func updateStreaks(userID, clubID uint) error {
	var entry LeaderboardEntry

	err := db.DB.
		Where("user_id = ? AND club_id = ?", userID, clubID).
		First(&entry).Error
	if err != nil {
		return err
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	switch {
	case entry.LastCheckedIn == nil:
		entry.CurrentStreak = 1
		if entry.LongestStreak < 1 {
			entry.LongestStreak = 1
		}
		entry.LastCheckedIn = &now

	default:
		lastDay := entry.LastCheckedIn.Truncate(24 * time.Hour)

		// Checked in yesterday → increment streak
		if lastDay.Equal(yesterday) {
			entry.CurrentStreak++
			if entry.CurrentStreak > entry.LongestStreak {
				entry.LongestStreak = entry.CurrentStreak
			}
			entry.LastCheckedIn = &now

			// Missed a day → reset
		} else if lastDay.Before(yesterday) {
			entry.CurrentStreak = 1
			if entry.LongestStreak < 1 {
				entry.LongestStreak = 1
			}
			entry.LastCheckedIn = &now

			// Same day → do nothing
		} else {
			return nil
		}
	}

	return db.DB.Save(&entry).Error
}

func UpdateLeaderboardScore(userID, clubID uint, delta int) error {
	if err := updateStreaks(userID, clubID); err != nil {
		return err
	}

	return db.DB.
		Model(&LeaderboardEntry{}).
		Where("user_id = ? AND club_id = ?", userID, clubID).
		UpdateColumns(map[string]interface{}{
			"score":          gorm.Expr("score + ?", delta),
			"last_checkedin": time.Now(),
		}).Error
}

func GetLeaderboardForClub(clubID uint, limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry

	err := db.DB.
		Where("club_id = ?", clubID).
		Order("score DESC, last_checkedin ASC").
		Limit(limit).
		Find(&entries).Error

	return entries, err
}

func GetLeaderboardEntryForUser(userID, clubID uint) (*LeaderboardEntry, error) {
	var entry LeaderboardEntry

	err := db.DB.
		Where("user_id = ? AND club_id = ?", userID, clubID).
		First(&entry).Error

	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func GetUserRankInClub(userID, clubID uint) (int, error) {
	var rank int

	query := `
	SELECT COUNT(*) + 1
	FROM leaderboard l
	JOIN leaderboard me
	  ON me.user_id = ? AND me.club_id = ?
	WHERE l.club_id = ?
	  AND (
	    l.score > me.score OR
	    (l.score = me.score AND l.last_checkedin > me.last_checkedin)
	  )
	`

	err := db.DB.Raw(query, userID, clubID, clubID).Scan(&rank).Error
	return rank, err
}

func CalculatePercentile(userID, clubID uint) (string, error) {
	var totalMembers int64

	if err := db.DB.
		Table("members").
		Where("club_id = ?", clubID).
		Count(&totalMembers).Error; err != nil {
		return "N/A", err
	}

	if totalMembers <= 1 {
		return "Top 100%", nil
	}

	rank, err := GetUserRankInClub(userID, clubID)
	if err != nil {
		return "N/A", err
	}

	percentage := (float64(rank) / float64(totalMembers)) * 100

	if percentage <= 1 {
		return "Top 1%", nil
	}

	return fmt.Sprintf("Top %.0f%%", percentage), nil
}

func GetWeeklyActivity(clubID, userID uint) (map[string]int, error) {
	rows, err := db.DB.Raw(`
		SELECT DATE(timestamp) as day, SUM(points)
		FROM activity_log
		WHERE club_id = ? AND user_id = ?
		  AND timestamp >= NOW() - INTERVAL '7 days'
		GROUP BY day
	`, clubID, userID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activity := make(map[string]int)
	for rows.Next() {
		var day string
		var points int
		rows.Scan(&day, &points)
		activity[day] = points
	}

	return activity, nil
}

func GetClubLeaderID(clubID uint) uint {
	var leaderID uint

	db.DB.
		Model(&LeaderboardEntry{}).
		Select("user_id").
		Where("club_id = ?", clubID).
		Order("score DESC").
		Limit(1).
		Scan(&leaderID)

	return leaderID
}
