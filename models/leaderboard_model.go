package models

import (
	"errors"
	"klubRanks/db"
	"time"
)

type LeaderboardEntry struct {
	ID            int64      `db:"id" json:"id"`
	UserID        int64      `db:"userid" json:"user_id"`
	ClubID        int64      `db:"clubid" json:"club_id"`
	Score         int        `db:"score" json:"score"`
	LastCheckedIn *time.Time `db:"last_checkedin" json:"last_checkedin,omitempty"`
}

func AddUserToLeaderboard(userID, clubID int64) error {
	query := `
	INSERT INTO leaderboard (userid, clubid, score, last_checkedin)
	VALUES (?, ?, ?, ?)
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
		time.Now(),
	)

	return err
}

func SetLeaderboardScore(userID, clubID int64, score int) error {
	query := `
	UPDATE leaderboard
	SET score = ?, last_checkedin = ?
	WHERE userid = ? AND clubid = ?
	`

	_, err := db.DB.Exec(
		query,
		score,
		time.Now(),
		userID,
		clubID,
	)

	return err
}

func UpdateLeaderboardScore(userID, clubID int64, delta int) error {
	query := `
	UPDATE leaderboard
	SET score = score + ?, last_checkedin = ?
	WHERE userid = ? AND clubid = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		delta,
		time.Now(),
		userID,
		clubID,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("leaderboard entry not found")
	}

	return nil
}

func GetLeaderboardForClub(clubID int64, limit int) ([]LeaderboardEntry, error) {
	query := `
	SELECT id, userid, clubid, score, last_checkedin
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
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}
