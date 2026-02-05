package api

import (
	"net/http"
	"time"

	"memory-quest-backend/internal/db"
	"memory-quest-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type AnswerPayload struct {
	IsCorrect bool `json:"is_correct"`
}

func SubmitQuizAnswer(c *gin.Context) {
	resultID := c.Param("resultId")
	userID := c.MustGet("user_id").(uint)

	var payload AnswerPayload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	var result models.QuizResult
	if err := db.DB.First(&result, resultID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz result not found"})
		return
	}

	// Only owner can submit answers
	if result.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your attempt"})
		return
	}

	// If already completed, reject
	if result.CompletedAt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quiz already completed"})
		return
	}

	now := time.Now()

	baseline := result.LastActivityAt
	if baseline == nil {
		// fall back to created_at (attempt start)
		t := result.CreatedAt
		baseline = &t
	}

	// Auto-finish if inactive for 1 hour
	if now.Sub(*baseline) > time.Hour {
		result.CompletedAt = &now
		db.DB.Save(&result)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Quiz attempt expired (inactive for over 1 hour)",
		})
		return
	}

	result.LastActivityAt = &now

	// Basic progress update
	result.AnsweredCount++
	if payload.IsCorrect {
		result.Score++
	}

	// If finished, mark completed
	if result.AnsweredCount >= result.TotalQuestions {
		result.CompletedAt = &now
	}

	if err := db.DB.Save(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quiz result"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quiz_result_id":  result.ID,
		"score":           result.Score,
		"answered_count":  result.AnsweredCount,
		"total_questions": result.TotalQuestions,
		"completed":       result.CompletedAt != nil,
	})
}
