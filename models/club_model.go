package models

import (
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"fmt"
	"strings"
	"time"

	"klubRanks/db"
	"klubRanks/logger"

	"gorm.io/gorm"
)

type Club struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	IsPrivate   bool      `json:"is_private"`
	Code        string    `gorm:"not null" json:"code"`
	Name        string    `gorm:"not null" json:"name"`
	Description *string   `json:"description,omitempty"`
	Action      string    `gorm:"not null" json:"action"`
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
		c.GenerateCode()

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

func (c *Club) GenerateCode() {
	// Use nanosecond timestamp for uniqueness
	payload := fmt.Sprintf("%d", c.CreatedAt.UnixNano())

	hash := sha256.Sum256([]byte(payload))

	c.Code = strings.ToUpper(
		base32.StdEncoding.WithPadding(base32.NoPadding).
			EncodeToString(hash[:5]),
	)

	logger.LogDebug("Generated code", c.Code, "at", c.CreatedAt)
}

func getClubByCode(clubCode string) (*Club, error) {
	var club Club

	err := db.DB.
		Where("code = ?", clubCode).
		First(&club).Error
	if err != nil {
		return nil, err
	}
	return &club, nil
}

func (c *Club) Update() error {
	return db.DB.
		Model(&Club{}).
		Where("id = ?", c.ID).
		Updates(map[string]interface{}{
			"name":        c.Name,
			"description": c.Description,
			"is_private":  c.IsPrivate,
			"action":      c.Action,
		}).Error
}

func getClubByID(clubID uint) (*Club, error) {
	var club Club

	err := db.DB.
		First(&club, clubID).Error
	if err != nil {
		return nil, err
	}
	return &club, nil
}

// Public wrapper for service layer
func GetClub(clubID uint) (*Club, error) {
	return getClubByID(clubID)
}

func AddMember(userID uint, clubCode string, role string) error {
	club, err := getClubByCode(clubCode)
	if isMember, _ := IsUserMemberOfClub(userID, club.ID); isMember {
		return errors.New("user is already a member of the club")
	}
	if err != nil {
		return errors.New("club not found")
	}
	member := Member{
		UserID:   userID,
		ClubID:   club.ID,
		Role:     role,
		JoinedAt: time.Now(),
	}

	if err := db.DB.Create(&member).Error; err != nil {
		return err
	}

	AddActivityLog(userID, club.ID, 0, ActionJoin)

	return AddUserToLeaderboard(userID, club.ID)
}

func IsUserMemberOfClub(userID, clubID uint) (bool, error) {
	var count int64

	err := db.DB.
		Model(&Member{}).
		Where("user_id = ? AND club_id = ?", userID, clubID).
		Count(&count).Error

	return count > 0, err
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
	AddActivityLog(userID, clubID, 0, ActionLeave)
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
