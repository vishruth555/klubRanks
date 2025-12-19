package models

import (
	"time"

	"klubRanks/db"
)

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"message_id"`
	ClubID    uint      `gorm:"not null;index" json:"club_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Timestamp time.Time `gorm:"not null;index" json:"timestamp"`
	Message   string    `gorm:"not null" json:"message"`
}

func AddMessage(m *Message) error {
	m.Timestamp = time.Now()
	return db.DB.Create(m).Error
}

func GetMessagesForClub(clubID uint, limit, offset int) ([]Message, error) {
	var messages []Message

	err := db.DB.
		Where("club_id = ?", clubID).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	return messages, err
}
