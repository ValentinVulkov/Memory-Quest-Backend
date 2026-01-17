package models

import "time"

type QuizResult struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"column:user_id;not null" json:"user_id"`
	QuizID         uint      `gorm:"column:quiz_id;not null" json:"quiz_id"`
	Score          int       `gorm:"default:0" json:"score"`
	TotalQuestions int       `gorm:"column:total_questions;default:0" json:"total_questions"`
	CompletedAt    time.Time `gorm:"column:completed_at" json:"completed_at"`
}
