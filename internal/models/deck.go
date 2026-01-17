package models

import "time"

type Deck struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"column:user_id;not null"`
	Title       string `gorm:"size:100;not null"`
	Description string
	IsPublic    bool      `gorm:"column:is_public"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`

	Cards []Card `gorm:"foreignKey:DeckID"`
}
