package api

import (
	"memory-quest-backend/internal/db"
	"memory-quest-backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsPublic    bool   `json:"is_public"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deck := models.Deck{
		UserID:      userID,
		Title:       body.Title,
		Description: body.Description,
		IsPublic:    body.IsPublic,
	}

	if err := db.DB.Create(&deck).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create deck"})
		return
	}

	c.JSON(201, deck)
}

func GetAllDecks(c *gin.Context) {
	var decks []models.Deck

	db.DB.Where("is_public = ?", true).Find(&decks)

	c.JSON(200, decks)
}

func GetDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id := c.Param("id")

	var deck models.Deck
	if err := db.DB.Preload("Cards").First(&deck, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Deck not found"})
		return
	}

	// Privacy: only owner can view private decks
	if !deck.IsPublic && deck.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(200, deck)
}

func UpdateDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id := c.Param("id")

	var deck models.Deck
	if err := db.DB.First(&deck, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Deck not found"})
		return
	}

	if deck.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsPublic    bool   `json:"is_public"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deck.Title = body.Title
	deck.Description = body.Description
	deck.IsPublic = body.IsPublic

	if err := db.DB.Save(&deck).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update deck"})
		return
	}
	c.JSON(200, deck)
}

func DeleteDeck(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)
	id := c.Param("id")

	var deck models.Deck

	if err := db.DB.Preload("Cards").First(&deck, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Deck not found"})
		return
	}

	if deck.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	db.DB.Delete(&deck)
	c.JSON(200, gin.H{"success": "Successfully deleted deck"})
}

func GetMyDecks(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var decks []models.Deck
	db.DB.Where("user_id = ?", userID).Find(&decks)

	c.JSON(200, decks)
}

func GetPublicDecks(c *gin.Context) {
	var decks []models.Deck

	if err := db.DB.Where("is_public = ?", true).Find(&decks).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch public decks"})
		return
	}

	c.JSON(200, decks)
}

func GetPublicDeck(c *gin.Context) {
	id := c.Param("id")

	var deck models.Deck
	if err := db.DB.First(&deck, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Deck not found"})
		return
	}

	if !deck.IsPublic {
		c.JSON(403, gin.H{"error": "Deck is private"})
		return
	}

	c.JSON(200, deck)
}
