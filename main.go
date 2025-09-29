package main

import (
	"log"
	"os"

	"devops-autopilot/routes"
	"devops-autopilot/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize OpenAI client
	utils.InitOpenAI()

	// Create Gin router
	r := gin.Default()

	// Setup all routes
	routes.SetupRoutes(r)

	// Get port from environment or default to 5000
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("🚀 Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
