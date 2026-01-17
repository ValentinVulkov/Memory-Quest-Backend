package api

import (
	"memory-quest-backend/internal/db"
	"memory-quest-backend/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateQuiz(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var body struct {
		DeckID uint   `json:"deck_id" binding:"required"`
		Title  string `json:"title"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var deck models.Deck
	if err := db.DB.Preload("Cards").First(&deck, body.DeckID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Deck not found"})
		return
	}

	if len(deck.Cards) == 0 {
		c.JSON(400, gin.H{"error": "Deck has no cards"})
		return
	}

	title := body.Title
	if title == "" {
		title = "Quiz: " + deck.Title
	}

	quiz := models.Quiz{
		UserID: userID,
		DeckID: deck.ID,
		Title:  title,
	}

	if err := db.DB.Create(&quiz).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create quiz"})
		return
	}

	var questions []models.QuizQuestion
	for _, card := range deck.Cards {
		question := models.QuizQuestion{
			QuizID:        quiz.ID,
			CardID:        card.ID,
			QuestionText:  card.Question,
			CorrectAnswer: card.Answer,
		}
		questions = append(questions, question)
	}

	if err := db.DB.Create(&questions).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create quiz question"})
		return
	}

	quiz.Questions = questions

	c.JSON(201, quiz)

}

func GetQuizzes(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var quizzes []models.Quiz
	db.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&quizzes)

	c.JSON(200, quizzes)
}

func GetQuiz(c *gin.Context) {
	quizID := c.Param("id")

	var quiz models.Quiz
	if err := db.DB.Preload("Questions").First(&quiz, quizID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Quiz not found"})
		return
	}

	c.JSON(200, quiz)
}

func DeleteQuiz(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	quizID := c.Param("id")

	var quiz models.Quiz
	if err := db.DB.First(&quiz, quizID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Quiz not found"})
		return
	}

	if quiz.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	db.DB.Where("quiz_id = ?", quizID).Delete(&models.QuizQuestion{})

	if err := db.DB.Delete(&quiz).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete quiz"})
		return
	}

	c.JSON(200, gin.H{"message": "Quiz deleted successfully"})

}

func PlayQuiz(c *gin.Context) {
	quizID := c.Param("id")

	var quiz models.Quiz
	if err := db.DB.Preload("Questions").First(&quiz, quizID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Quiz not found"})
		return
	}

	type QuestionResponse struct {
		ID           uint   `json:"id"`
		QuestionText string `json:"question"`
	}

	var questions []QuestionResponse
	for _, q := range quiz.Questions {
		questions = append(questions, QuestionResponse{
			ID:           q.ID,
			QuestionText: q.QuestionText,
		})
	}

	c.JSON(200, gin.H{
		"quiz_id":   quiz.ID,
		"title":     quiz.Title,
		"questions": questions,
	})
}
