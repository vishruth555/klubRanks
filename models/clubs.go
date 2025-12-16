package models

import (
	"klubRanks/db"
	"time"
)

type Club struct {
	ID          int64     `db:"id" json:"id"`
	CreatedBy   int64     `db:"created_by" json:"created_by"`
	IsPrivate   bool      `db:"is_private" json:"is_private"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type Member struct {
	ID       int64     `db:"id" json:"id"`
	UserID   int64     `db:"userid" json:"user_id"`
	ClubID   int64     `db:"clubid" json:"club_id"`
	Role     string    `db:"role" json:"role"`
	JoinedAt time.Time `db:"joined_at" json:"joined_at"`
}

func (c *Club) Save() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	query := `
	INSERT INTO clubs (created_by, is_private, name, description, created_at)
	VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	now := time.Now()

	result, err := stmt.Exec(
		c.CreatedBy,
		c.IsPrivate,
		c.Name,
		c.Description,
		now,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	c.ID = id
	c.CreatedAt = now

	// Add creator as admin member
	memberQuery := `
	INSERT INTO members (userid, clubid, role, joined_at)
	VALUES (?, ?, ?, ?)
	`

	_, err = tx.Exec(
		memberQuery,
		c.CreatedBy,
		c.ID,
		"admin",
		now,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *Club) Update() error {
	query := `
	UPDATE clubs
	SET name = ?, description = ?, is_private = ?
	WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		c.Name,
		c.Description,
		c.IsPrivate,
		c.ID,
	)

	return err
}

func AddMember(userID, clubID int, role string) error {
	query := `
	INSERT INTO members (userid, clubid, role, joined_at)
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
		role,
		time.Now(),
	)

	return err
}

func GetClubsForUser(userID int) ([]Club, error) {
	query := `
	SELECT 
		c.id,
		c.created_by,
		c.is_private,
		c.name,
		c.description,
		c.created_at
	FROM clubs c
	INNER JOIN members m ON m.clubid = c.id
	WHERE m.userid = ?
	ORDER BY c.created_at DESC
	`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clubs []Club

	for rows.Next() {
		var club Club
		err := rows.Scan(
			&club.ID,
			&club.CreatedBy,
			&club.IsPrivate,
			&club.Name,
			&club.Description,
			&club.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		clubs = append(clubs, club)
	}

	return clubs, nil
}

func RemoveMember(userID, clubID int) error {
	query := `
	DELETE FROM members
	WHERE userid = ? AND clubid = ?
	`
	_, err := db.DB.Exec(query, userID, clubID)
	return err
}

func GetClubMembers(clubID int) ([]Member, error) {
	query := `
	SELECT id, userid, clubid, role, joined_at
	FROM members
	WHERE clubid = ?
	`

	rows, err := db.DB.Query(query, clubID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member

	for rows.Next() {
		var m Member
		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.ClubID,
			&m.Role,
			&m.JoinedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	return members, nil
}
