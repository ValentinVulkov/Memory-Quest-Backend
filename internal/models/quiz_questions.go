package models

import "time"

// QuizQuestion represents a single question generated for a Quiz.
// Storing questions preserves the quiz snapshot even if the deck changes later.
type QuizQuestion struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	QuizID        uint   `gorm:"not null;index" json:"quiz_id"`
	CardID        uint   `gorm:"not null;index" json:"card_id"`
	QuestionText  string `gorm:"type:text" json:"question_text"`
	CorrectAnswer string `gorm:"type:text" json:"correct_answer"`
	// Optionally: store choices JSON or other metadata
	CreatedAt time.Time `json:"created_at"`
}
