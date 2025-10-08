package routes

import (
	"devops-autopilot/handlers"

	"github.com/gin-gonic/gin"
)

// SetupProvisionRoutes sets up all provision-related routes
func SetupProvisionRoutes(router *gin.RouterGroup) {
	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)
	
	// Terraform generation endpoint (OpenAI)
	router.POST("/terraform", handlers.GenerateTerraform)
	
	// Terraform generation endpoint (GitHub Copilot)
	router.POST("/terraform-copilot", handlers.GenerateTerraformWithCopilot)
	
	// Terraform validation endpoint  
	router.POST("/validate", handlers.ValidateTerraform)
}

// SetupRoutes sets up all application routes
func SetupRoutes(r *gin.Engine) {
	// API group
	api := r.Group("/api/provision")
	SetupProvisionRoutes(api)
	
	// Future route groups can be added here
	// v2 := r.Group("/api/v2")
	// auth := r.Group("/auth")
}