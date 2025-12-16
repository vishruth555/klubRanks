package models

import (
	"klubRanks/db"
	"time"
)

type Message struct {
	MessageID int64     `db:"messageid" json:"message_id"`
	ClubID    int64     `db:"clubid" json:"club_id"`
	UserID    int64     `db:"userid" json:"user_id"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	Message   string    `db:"message" json:"message"`
}

func AddMessage(m *Message) error {
	query := `
	INSERT INTO messages (clubid, userid, timestamp, message)
	VALUES (?, ?, ?, ?)
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()

	result, err := stmt.Exec(
		m.ClubID,
		m.UserID,
		now,
		m.Message,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	m.MessageID = id
	m.Timestamp = now

	return nil
}

func GetMessagesForClub(clubID int64, limit int, offset int) ([]Message, error) {
	query := `
	SELECT messageid, clubid, userid, timestamp, message
	FROM messages
	WHERE clubid = ?
	ORDER BY timestamp DESC
	LIMIT ? OFFSET ?
	`

	rows, err := db.DB.Query(query, clubID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var m Message
		err := rows.Scan(
			&m.MessageID,
			&m.ClubID,
			&m.UserID,
			&m.Timestamp,
			&m.Message,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}
