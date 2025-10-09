package handlers

import (
	"net/http"
	"strings"

	"devops-autopilot/models"
	"devops-autopilot/services"
	"devops-autopilot/utils"

	"github.com/gin-gonic/gin"
)

var terraformService = services.NewTerraformService()

// HealthCheck handles the health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{
		Status: "Service is healthy test",
	})
}

// ValidateTerraform handles terraform code validation
func ValidateTerraform(c *gin.Context) {
	var req models.ValidationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.TerraformCode) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "terraformCode field cannot be empty",
		})
		return
	}

	// Validate the provided Terraform code
	validation, err := utils.ValidateTerraformCode(req.TerraformCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to validate terraform code",
			"details": err.Error(),
		})
		return
	}

	// Return validation results
	statusCode := http.StatusOK
	if !validation.IsValid {
		statusCode = http.StatusUnprocessableEntity // 422 - validation failed
	}

	c.JSON(statusCode, models.ValidationResponse{
		Message:    "Terraform validation completed",
		Validation: validation,
	})
}

// GenerateTerraform handles terraform code generation
func GenerateTerraform(c *gin.Context) {
	var req models.TerraformRequest

	// Validate JSON input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.Resource) == "" || strings.TrimSpace(req.Specs) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Resource and specs fields cannot be empty",
		})
		return
	}

	// Generate and validate terraform code
	cleanedCode, validation, err := terraformService.GenerateAndValidate(req.Resource, req.Specs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Save file only if validation passes
	if validation.IsValid {
		_, err := terraformService.SaveTerraformFile(cleanedCode, req.Resource, "openai")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to save terraform file",
				"details": err.Error(),
			})
			return
		}
	}

	// Determine response status and message based on validation
	statusCode := http.StatusOK
	message := "Terraform code generated successfully"

	if !validation.IsValid {
		statusCode = http.StatusCreated // 201 - generated but has validation errors
		message = "Terraform code generated with validation errors"
	}

	// Success response with validation results
	c.JSON(statusCode, models.TerraformResponse{
		Message:       message,
		TerraformCode: cleanedCode,
		Validation:    validation,
	})
}

// GenerateTerraformWithCopilot handles terraform code generation using GitHub Copilot
func GenerateTerraformWithCopilot(c *gin.Context) {
	var req models.TerraformRequest

	// Validate JSON input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.Resource) == "" || strings.TrimSpace(req.Specs) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Resource and specs fields cannot be empty",
		})
		return
	}

	// Generate and validate terraform code using GitHub Copilot
	cleanedCode, validation, err := terraformService.GenerateAndValidateWithCopilot(req.Resource, req.Specs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Save file only if validation passes
	if validation.IsValid {
		_, err := terraformService.SaveTerraformFile(cleanedCode, req.Resource, "copilot")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to save terraform file",
				"details": err.Error(),
			})
			return
		}
	}

	// Determine response status and message based on validation
	statusCode := http.StatusOK
	message := "Terraform code generated successfully using GitHub Copilot"

	if !validation.IsValid {
		statusCode = http.StatusCreated // 201 - generated but has validation errors
		message = "Terraform code generated using GitHub Copilot with validation errors"
	}

	// Success response with validation results
	c.JSON(statusCode, models.TerraformResponse{
		Message:       message,
		TerraformCode: cleanedCode,
		Validation:    validation,
	})
}
