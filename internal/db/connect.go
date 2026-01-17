package db

import (
	"log"
	"memory-quest-backend/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Replace with your actual MySQL credentials
	username := "root"   // Your MySQL username
	password := "Pas123" // Your MySQL password
	host := "127.0.0.1"
	port := "3306"

	// Connect without specifying database first
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	// Create database
	db.Exec("CREATE DATABASE IF NOT EXISTS memory_quest CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
	db.Exec("USE memory_quest")

	// Reconnect to the specific database
	dsn = username + ":" + password + "@tcp(" + host + ":" + port + ")/memory_quest?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to memory_quest database:", err)
	}

	// AutoMigrate all models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Deck{},
		&models.Card{},
		&models.Quiz{},
		&models.QuizQuestion{},
		&models.QuizResult{},
		&models.Leaderboard{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully!")
}
