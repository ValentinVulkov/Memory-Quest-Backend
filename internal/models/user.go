package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"size:50;not null;unique"`
	Email        string    `gorm:"size:100;not null;unique"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	Role         string    `gorm:"type:enum('user','admin');default:'user'"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}
