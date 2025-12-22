package models

import (
	"time"

	"klubRanks/db"
)

const (
	MessageTypeUser   = "user"   // normal chat
	MessageTypeSystem = "system" // join/leave/score update
)

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"message_id"`
	ClubID    uint      `gorm:"not null;index" json:"club_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"` // to preload user info
	Timestamp time.Time `gorm:"not null;index" json:"timestamp"`
	Type      string    `gorm:"not null;index" json:"type"`
	Message   string    `gorm:"not null" json:"message"`
	ReplyToID *uint     `json:"reply_to_id"`
	ReplyTo   *Message  `gorm:"foreignKey:ReplyToID" json:"reply_to,omitempty"`
}

func (m *Message) AddMessage() error {
	m.Timestamp = time.Now()
	return db.DB.Create(m).Error
}

func GetMessagesForClub(clubID uint, limit, offset int) ([]Message, error) {
	var messages []Message

	err := db.DB.
		Preload("User").
		Preload("ReplyTo").
		Preload("ReplyTo.User").
		Where("club_id = ?", clubID).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	return messages, err
}
