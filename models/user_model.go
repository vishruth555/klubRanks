package models

import (
	"errors"
	"klubRanks/db"
	"klubRanks/utils"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	AvatarID  string    `gorm:"default:default" json:"avatar_id"`
	CreatedAt time.Time `json:"created_at"`
}

func UpdateAvatar(userID uint, avatarID string) error {
	return db.DB.
		Model(&User{}).
		Where("id = ?", userID).
		Update("avatar_id", avatarID).
		Error
}

func GetUserByID(id uint) (*User, error) {
	var user User

	err := db.DB.
		Select("id", "username", "avatar_id", "created_at").
		First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) Save() error {
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hashedPassword
	u.CreatedAt = time.Now()

	return db.DB.Create(u).Error
}

func (u *User) ValidateCredentials() error {
	var user User

	err := db.DB.
		Where("username = ?", u.Username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid credentials")
		}
		return err
	}

	if !utils.CheckPasswordHash(u.Password, user.Password) {
		return errors.New("invalid credentials")
	}

	u.ID = user.ID
	u.AvatarID = user.AvatarID

	return nil
}
