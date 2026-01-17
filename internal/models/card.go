package models

import "time"

type Card struct {
	ID        uint      `gorm:"primaryKey"`
	DeckID    uint      `gorm:"column:deck_id;not null"`
	Question  string    `gorm:"not null"`
	Answer    string    `gorm:"not null"`
	ImageURL  string    `gorm:"column:image_url"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
