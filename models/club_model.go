package models

import (
	"time"

	"klubRanks/db"

	"gorm.io/gorm"
)

type Club struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	IsPrivate   bool      `json:"is_private"`
	Name        string    `gorm:"not null" json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`

	Members []Member `gorm:"foreignKey:ClubID"`
}

type Member struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `gorm:"not null" json:"user_id"`
	ClubID   uint      `gorm:"not null;index" json:"club_id"`
	Role     string    `gorm:"not null" json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

func (c *Club) Save() error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		c.CreatedAt = time.Now()

		if err := tx.Create(c).Error; err != nil {
			return err
		}

		member := Member{
			UserID:   c.CreatedBy,
			ClubID:   c.ID,
			Role:     "admin",
			JoinedAt: time.Now(),
		}

		if err := tx.Create(&member).Error; err != nil {
			return err
		}

		return nil
	})
}

func (c *Club) Update() error {
	return db.DB.
		Model(&Club{}).
		Where("id = ?", c.ID).
		Updates(map[string]interface{}{
			"name":        c.Name,
			"description": c.Description,
			"is_private":  c.IsPrivate,
		}).Error
}

func AddMember(userID, clubID uint, role string) error {
	member := Member{
		UserID:   userID,
		ClubID:   clubID,
		Role:     role,
		JoinedAt: time.Now(),
	}

	if err := db.DB.Create(&member).Error; err != nil {
		return err
	}

	return AddUserToLeaderboard(userID, clubID)
}

func GetClubsForUser(userID uint) ([]Club, error) {
	var clubs []Club

	err := db.DB.
		Joins("JOIN members ON members.club_id = clubs.id").
		Where("members.user_id = ?", userID).
		Order("clubs.created_at ASC").
		Find(&clubs).Error

	return clubs, err
}

func GetMemberCountForClub(clubID uint) (int64, error) {
	var count int64

	err := db.DB.
		Model(&Member{}).
		Where("club_id = ?", clubID).
		Count(&count).Error

	return count, err
}

func RemoveMember(userID, clubID uint) error {
	return db.DB.
		Where("user_id = ? AND club_id = ?", userID, clubID).
		Delete(&Member{}).
		Error
}

func GetClubMembers(clubID uint) ([]Member, error) {
	var members []Member

	err := db.DB.
		Where("club_id = ?", clubID).
		Find(&members).Error

	return members, err
}
