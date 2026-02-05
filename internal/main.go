package main

import (
	"log"
	"os"
	"time"

	"memory-quest-backend/internal/api"
	"memory-quest-backend/internal/db"
	"memory-quest-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load env from either repo root or /internal working dir
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	// Hard safety: don't allow empty secret (prevents "everything broke" again)
	if os.Getenv("JWT_SECRET") == "" {
		// Match your .env default
		os.Setenv("JWT_SECRET", "supersecretjwt")
		log.Println("⚠️ JWT_SECRET was missing; using dev default (supersecretjwt)")
	}

	db.Connect()

	r := gin.Default()

	// CORS for Vite dev server
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Auth
	r.POST("/api/register", api.Register)
	r.POST("/api/login", api.Login)

	// Public decks (no auth required)
	r.GET("/api/decks/public", api.GetAllDecks)
	r.GET("/api/decks/public/:id", api.GetPublicDeck)
	r.GET("/api/decks/public/:id/cards", api.GetPublicCards)

	auth := r.Group("/api")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/profile", func(c *gin.Context) {
			userID := c.MustGet("user_id").(uint)
			c.JSON(200, gin.H{"user_id": userID})
		})

		// Decks
		auth.GET("/decks", api.GetMyDecks)
		auth.POST("/decks", api.CreateDeck)
		auth.GET("/decks/:id", api.GetDeck)
		auth.PUT("/decks/:id", api.UpdateDeck)
		auth.DELETE("/decks/:id", api.DeleteDeck)

		// Cards
		auth.GET("/decks/:id/cards", api.GetCards)
		auth.POST("/decks/:id/cards", api.CreateCard)
		auth.GET("/decks/:id/cards/:cardId", api.GetCard)
		auth.PUT("/decks/:id/cards/:cardId", api.UpdateCard)
		auth.DELETE("/decks/:id/cards/:cardId", api.DeleteCard)

		//Quiz
		auth.POST("/decks/:id/quiz/start", api.StartQuiz)
		auth.POST("/quizzes/:resultId/answer", api.SubmitQuizAnswer)
	}
	//Continuesly check time.
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()
			expiry := now.Add(-1 * time.Hour)

			// Expire attempts with last_activity_at older than 1 hour
			db.DB.Exec(`
			UPDATE quiz_results
			SET completed_at = ?
			WHERE completed_at IS NULL
			  AND last_activity_at IS NOT NULL
			  AND last_activity_at < ?
		`, now, expiry)

			// Expire attempts that never had activity, based on created_at
			db.DB.Exec(`
			UPDATE quiz_results
			SET completed_at = ?
			WHERE completed_at IS NULL
			  AND last_activity_at IS NULL
			  AND created_at < ?
		`, now, expiry)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
