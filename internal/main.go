package main

import (
	"memory-quest-backend/internal/api"
	"memory-quest-backend/internal/db"
	"memory-quest-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	r := gin.Default()

	// Auth
	r.POST("/api/register", api.Register)
	r.POST("/api/login", api.Login)

	r.GET("/api/decks", api.GetAllDecks)
	r.GET("/api/decks/:id", api.GetDeck)
	r.GET("/api/decks/:id/cards", api.GetCards)

	auth := r.Group("/api")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/profile", func(c *gin.Context) {
			userID := c.MustGet("user_id").(uint)
			c.JSON(200, gin.H{"user_id": userID})
		})

		auth.POST("/decks", api.CreateDeck)
		auth.GET("/my-decks", api.GetMyDecks)
		auth.PUT("/decks/:id", api.UpdateDeck)
		auth.DELETE("/decks/:id", api.DeleteDeck)

		auth.POST("/decks/:id/cards", api.CreateCard)
		auth.GET("/decks/:id/cards/:cardId", api.GetCard)
		auth.PUT("/decks/:id/cards/:cardId", api.UpdateCard)
		auth.DELETE("/decks/:id/cards/:cardId", api.DeleteCard)
	}

	r.Run(":8080")
}
