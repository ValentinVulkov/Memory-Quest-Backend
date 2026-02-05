package models

import "time"

type QuizResult struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	QuizID         uint       `json:"quiz_id"`
	UserID         uint       `json:"user_id"`
	Score          int        `json:"score"`
	TotalQuestions int        `json:"total_questions"`
	CompletedAt    *time.Time `json:"completed_at"` // ✅ nullable
	CreatedAt      time.Time  `json:"created_at"`
	AnsweredCount  int        `json:"answered_count" gorm:"default:0"`
	LastActivityAt *time.Time `json:"last_activity_at"`
}
