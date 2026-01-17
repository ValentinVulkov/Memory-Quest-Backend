package models

import "time"

// Leaderboard is an optional cached table for fast queries.
// You can also compute leaderboard dynamically from quiz_results.
type Leaderboard struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	UserID                uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	TotalQuizzesCompleted int       `gorm:"default:0" json:"total_quizzes_completed"`
	HighestScore          int       `gorm:"default:0" json:"highest_score"`
	AverageAccuracy       float64   `gorm:"type:decimal(6,2);default:0" json:"average_accuracy"`
	LastUpdated           time.Time `gorm:"autoUpdateTime" json:"last_updated"`
}
