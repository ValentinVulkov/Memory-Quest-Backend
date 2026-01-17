package api

import (
	"memory-quest-backend/internal/auth"
	"memory-quest-backend/internal/db"
	"memory-quest-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid Input"})
		return
	}

	var existingUser models.User
	if err := db.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(400, gin.H{"error": "User with this email already exists"})
		return
	}

	if err := db.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(400, gin.H{"errpr": "User with this email already exists"})
		return
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully"})

}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid Input"})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or passowrd"})
		return
	}

	if !auth.CheckPasswordHash(input.Password, user.PasswordHash) {
		c.JSON(401, gin.H{"error": " email or password"})
		return
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Logged in successfully",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})

}
