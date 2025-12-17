package models

import (
	"errors"
	"klubRanks/db"
	"klubRanks/utils"
	"time"
)

type User struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"-"`
	AvatarID  string    `db:"avatar_id" json:"avatar_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func GetUserByID(id int) (*User, error) {
	query := `
	SELECT id, username, avatar_id, created_at
	FROM users
	WHERE id = ?
	`

	row := db.DB.QueryRow(query, id)

	var user User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.AvatarID,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) Save() error {
	query := `
	INSERT INTO users (username, password, avatar_id, created_at)
	VALUES (?, ?, ?, ?)
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	now := time.Now()

	result, err := stmt.Exec(
		u.Username,
		hashedPassword,
		u.AvatarID,
		now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id
	u.CreatedAt = now

	return nil
}

func (u *User) ValidateCredentials() error {
	query := `
	SELECT id, password
	FROM users
	WHERE username = ?
	`

	row := db.DB.QueryRow(query, u.Username)

	var hashedPassword string
	err := row.Scan(&u.ID, &hashedPassword)
	if err != nil {
		return errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(u.Password, hashedPassword) {
		return errors.New("invalid credentials")
	}

	return nil
}
