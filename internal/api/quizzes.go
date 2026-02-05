package api

import (
	"math/rand"
	"net/http"
	"time"

	"memory-quest-backend/internal/db"
	"memory-quest-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type QuizQuestionDTO struct {
	ID       uint     `json:"id"`
	CardID   uint     `json:"card_id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type StartQuizResponse struct {
	QuizID         uint              `json:"quiz_id"`
	QuizResultID   uint              `json:"quiz_result_id"`
	DeckID         uint              `json:"deck_id"`
	Title          string            `json:"title"`
	TotalQuestions int               `json:"total_questions"`
	Questions      []QuizQuestionDTO `json:"questions"`
}

func StartQuiz(c *gin.Context) {
	deckID := c.Param("id")
	userID := c.MustGet("user_id").(uint)

	// Load deck
	var deck models.Deck
	if err := db.DB.First(&deck, deckID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Deck not found"})
		return
	}

	// Permission check
	if !deck.IsPublic && deck.UserID != userID {
		c.JSON(403, gin.H{"error": "Not your deck"})
		return
	}

	// Load cards
	var cards []models.Card
	if err := db.DB.Where("deck_id = ?", deck.ID).Find(&cards).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to load cards"})
		return
	}

	if len(cards) < 4 {
		c.JSON(400, gin.H{"error": "Need at least 4 cards to start quiz"})
		return
	}

	// Build answer pool
	allAnswers := make([]string, 0, len(cards))
	for _, card := range cards {
		allAnswers = append(allAnswers, card.Answer)
	}

	// Create quiz
	quiz := models.Quiz{
		DeckID:    deck.ID,
		Title:     "Quiz: " + deck.Title,
		CreatedAt: time.Now(),
	}
	if err := db.DB.Create(&quiz).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create quiz"})
		return
	}

	// Create result entry immediately
	now := time.Now()

	result := models.QuizResult{
		QuizID:         quiz.ID,
		UserID:         userID,
		Score:          0,
		TotalQuestions: len(cards),
		AnsweredCount:  0,
		LastActivityAt: &now,
	}
	if err := db.DB.Create(&result).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create quiz result"})
		return
	}

	db.DB.Create(&result)

	questions := make([]QuizQuestionDTO, 0, len(cards))

	rand.Seed(time.Now().UnixNano())

	for _, card := range cards {
		opts := buildOptions(allAnswers, card.Answer)

		q := models.QuizQuestion{
			QuizID:        quiz.ID,
			CardID:        card.ID,
			QuestionText:  card.Question,
			CorrectAnswer: card.Answer,
		}
		db.DB.Create(&q)

		questions = append(questions, QuizQuestionDTO{
			ID:       q.ID,
			CardID:   card.ID,
			Question: card.Question,
			Options:  opts,
		})
	}

	c.JSON(http.StatusCreated, StartQuizResponse{
		QuizID:         quiz.ID,
		QuizResultID:   result.ID,
		DeckID:         deck.ID,
		Title:          quiz.Title,
		TotalQuestions: len(questions),
		Questions:      questions,
	})
}

func buildOptions(pool []string, correct string) []string {
	opts := []string{correct}
	used := map[string]bool{correct: true}

	for len(opts) < 4 {
		candidate := pool[rand.Intn(len(pool))]
		if !used[candidate] {
			opts = append(opts, candidate)
			used[candidate] = true
		}
	}

	rand.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })
	return opts
}
