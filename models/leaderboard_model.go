package models

import (
	"fmt"
	"klubRanks/db"
	"time"
)

type LeaderboardEntry struct {
	ID            int64      `db:"id" json:"id"`
	UserID        int64      `db:"userid" json:"user_id"`
	ClubID        int64      `db:"clubid" json:"club_id"`
	Score         int        `db:"score" json:"score"`
	CurrentStreak int        `db:"current_streak" json:"current_streak"`
	LongestStreak int        `db:"longest_streak" json:"longest_streak"`
	LastCheckedIn *time.Time `db:"last_checkedin" json:"last_checkedin,omitempty"`
}

func AddUserToLeaderboard(userID, clubID int64) error {
	query := `
	INSERT INTO leaderboard (userid, clubid, score)
	VALUES (?, ?, ?)
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		userID,
		clubID,
		0,
	)

	return err
}

func updateStreaks(userID, clubID int64) error {
	var lastCheckedIn *time.Time

	err := db.DB.QueryRow(
		`SELECT last_checkedin FROM leaderboard WHERE userid = ? AND clubid = ?`,
		userID, clubID,
	).Scan(&lastCheckedIn)

	if err != nil {
		return err
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	if lastCheckedIn == nil {
		_, err = db.DB.Exec(`
			UPDATE leaderboard
			SET
				current_streak = 1,
				longest_streak = CASE
					WHEN longest_streak < 1 THEN 1
					ELSE longest_streak
				END,
				last_checkedin = ?
			WHERE userid = ? AND clubid = ?
		`, now, userID, clubID)
		return err
	}

	lastDay := lastCheckedIn.Truncate(24 * time.Hour)

	switch {

	// Checked in yesterday → increment streak
	case lastDay.Equal(yesterday):
		_, err = db.DB.Exec(`
			UPDATE leaderboard
			SET current_streak = current_streak + 1,
			    longest_streak = CASE
			        WHEN current_streak + 1 > longest_streak
			        THEN current_streak + 1
			        ELSE longest_streak
			    END,
			    last_checkedin = ?
			WHERE userid = ? AND clubid = ?
		`, now, userID, clubID)

	// Missed a day → reset streak
	case lastDay.Before(yesterday):
		_, err = db.DB.Exec(`
			UPDATE leaderboard
			SET
				current_streak = 1,
				longest_streak = CASE
					WHEN longest_streak < 1 THEN 1
					ELSE longest_streak
				END,
				last_checkedin = ?
			WHERE userid = ? AND clubid = ?
		`, now, userID, clubID)

	// Same day → do nothing
	default:
		return nil
	}

	return err
}

func UpdateLeaderboardScore(userID, clubID int64, delta int) error {
	err := updateStreaks(userID, clubID)
	if err != nil {
		return err
	}

	query := `
	UPDATE leaderboard
	SET score = score + ?, last_checkedin = ?
	WHERE userid = ? AND clubid = ?
	`

	_, err = db.DB.Exec(
		query,
		delta,
		time.Now(),
		userID,
		clubID,
	)

	if err != nil {
		return err
	}

	return err
}

func GetLeaderboardForClub(clubID int64, limit int) ([]LeaderboardEntry, error) {
	query := `
	SELECT id, userid, clubid, score, last_checkedin, current_streak
	FROM leaderboard
	WHERE clubid = ?
	ORDER BY score DESC, last_checkedin ASC
	LIMIT ?
	`

	rows, err := db.DB.Query(query, clubID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []LeaderboardEntry

	for rows.Next() {
		var e LeaderboardEntry
		err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.ClubID,
			&e.Score,
			&e.LastCheckedIn,
			&e.CurrentStreak,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func GetLeaderboardEntryForUser(userID, clubID int64) (*LeaderboardEntry, error) {
	query := `
	SELECT
		id,
		userid,
		clubid,
		score,
		current_streak,
		longest_streak,
		last_checkedin
	FROM leaderboard
	WHERE userid = ? AND clubid = ?
	LIMIT 1
	`

	var e LeaderboardEntry
	err := db.DB.QueryRow(query, userID, clubID).Scan(
		&e.ID,
		&e.UserID,
		&e.ClubID,
		&e.Score,
		&e.CurrentStreak,
		&e.LongestStreak,
		&e.LastCheckedIn,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func GetUserRankInClub(userID, clubID int64) (int, error) {
	query := `
	SELECT COUNT(*) + 1
	FROM leaderboard l
	JOIN leaderboard me
	  ON me.userid = ? AND me.clubid = ?
	WHERE l.clubid = ?
	  AND (
	    l.score > me.score OR
	    (l.score = me.score AND l.last_checkedin < me.last_checkedin)
	  )
	`

	var rank int
	err := db.DB.QueryRow(query, userID, clubID, clubID).Scan(&rank)
	if err != nil {
		return 0, err
	}

	return rank, nil
}

func CalculatePercentile(userID, clubID int64) (string, error) {
	var totalMembers int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM members WHERE clubid = ?`, clubID).Scan(&totalMembers)
	if err != nil {
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

func GetWeeklyActivity(clubID, userID int64) (map[string]int, error) {
	query := `
		SELECT strftime('%Y-%m-%d', timestamp) as day, SUM(points)
		FROM activity_log
		WHERE clubid = ? AND userid = ? AND timestamp >= datetime('now', '-7 days')
		GROUP BY day
	`
	rows, err := db.DB.Query(query, clubID, userID)
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

func GetClubLeaderID(clubID int64) int64 {
	var leaderID int64
	query := `SELECT userid FROM leaderboard WHERE clubid = ? ORDER BY score DESC LIMIT 1`
	err := db.DB.QueryRow(query, clubID).Scan(&leaderID)
	if err != nil {
		return 0
	}
	return leaderID
}