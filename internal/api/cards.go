package api

import (
	"memory-quest-backend/internal/db"
	"memory-quest-backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCard(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)
	deckID := c.Param("id")

	var deck models.Deck

	if err := db.DB.First(&deck, deckID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Deck not found"})
		return
	}

	if deck.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	var body struct {
		Question string `json:"question" binding:"required"`
		Answer   string `json:"answer" binding:"required"`
		ImageURL string `json:"image_url"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	card := models.Card{
		DeckID:   deck.ID,
		Question: body.Question,
		Answer:   body.Answer,
		ImageURL: body.ImageURL,
	}

	if err := db.DB.Create(&card).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create card"})
		return
	}

	c.JSON(201, card)

}

func GetCards(c *gin.Context) {
	deckID := c.Param("id")

	//Verification
	var deck models.Deck
	if err := db.DB.First(&deck, deckID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Deck not found"})
		return
	}

	var cards []models.Card
	db.DB.Where("deck_id = ?", deck.ID).Find(&cards)

	c.JSON(200, cards)
}

func GetCard(c *gin.Context) {
	deckID := c.Param("id")
	cardID := c.Param("cardId")

	var card models.Card
	if err := db.DB.Where("id = ? AND deck_id = ?", cardID, deckID).First(&card).Error; err != nil {
		c.JSON(404, gin.H{"error": "Card not found"})
		return
	}

	c.JSON(200, card)
}

func UpdateCard(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	deckID := c.Param("id")
	cardID := c.Param("cardId")

	var deck models.Deck
	if err := db.DB.First(&deck, deckID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Deck not found"})
		return
	}
	if deck.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
	}

	var card models.Card
	if err := db.DB.Where("id = ? AND deck_id = ?", cardID, deckID).First(&card).Error; err != nil {
		c.JSON(404, gin.H{"error": "Card not found"})
		return
	}

	var body struct {
		Question string `json:"question" binding:"required"`
		Answer   string `json:"answer" binding:"required"`
		ImageURL string `json:"image_url"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	card.Question = body.Question
	card.Answer = body.Answer
	card.ImageURL = body.ImageURL

	if err := db.DB.Save(&card).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update card"})
		return
	}
	c.JSON(200, card)
}

func DeleteCard(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	deckID := c.Param("id")
	cardID := c.Param("cardId")

	var deck models.Deck
	if err := db.DB.First(&deck, deckID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Deck not found"})
		return
	}

	if deck.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	var card models.Card
	if err := db.DB.Where("id = ? AND deck_id = ?", cardID, deckID).First(&card).Error; err != nil {
		c.JSON(404, gin.H{"error": "Card not found"})
		return
	}

	if err := db.DB.Delete(&card).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete card"})
		return
	}

	c.JSON(200, gin.H{"message": "Card deleted successfully"})
}

func GetPublicCards(c *gin.Context) {
	deckID := c.Param("id")

	var deck models.Deck
	if err := db.DB.First(&deck, deckID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Deck not found"})
		return
	}

	if !deck.IsPublic {
		c.JSON(403, gin.H{"error": "Deck is private"})
		return
	}

	var cards []models.Card
	db.DB.Where("deck_id = ?", deck.ID).Find(&cards)

	c.JSON(200, cards)
}
