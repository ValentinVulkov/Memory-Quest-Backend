package models

import "time"

// Quiz represents a quiz generated from a Deck. It can be considered a snapshot or template.

type Quiz struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id;not null" json:"user_id"`
	DeckID    uint      `gorm:"column:deck_id;not null" json:"deck_id"`
	Title     string    `gorm:"size:100" json:"title"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`

	Questions []QuizQuestion `gorm:"foreignkey:QuizID" json:"questions"`
}
