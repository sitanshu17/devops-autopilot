package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"devops-autopilot/utils"

	"github.com/gin-gonic/gin"
)

// TerraformRequest represents the request body for terraform generation
type TerraformRequest struct {
	Resource string `json:"resource" binding:"required"`
	Specs    string `json:"specs" binding:"required"`
}

// TerraformResponse represents the response for terraform generation
type TerraformResponse struct {
	Message       string `json:"message"`
	TerraformCode string `json:"terraformCode"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthCheck handles the health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status: "Service is healthy test",
	})
}

// GenerateTerraform handles terraform code generation
func GenerateTerraform(c *gin.Context) {
	var req TerraformRequest
	
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

	// Generate terraform code using OpenAI
	tfCode, err := utils.GenerateTerraformCode(req.Resource, req.Specs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate Terraform code",
			"details": err.Error(),
		})
		return
	}

	// Validate generated code is not empty
	if strings.TrimSpace(tfCode) == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Generated Terraform code is empty",
		})
		return
	}

	// Clean the code (remove markdown code block markers)
	cleanedCode, err := cleanTerraformCode(tfCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to clean Terraform code",
			"details": err.Error(),
		})
		return
	}

	// Ensure terraform directory exists
	terraformDir := "terraform"
	if err := os.MkdirAll(terraformDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create terraform directory",
			"details": err.Error(),
		})
		return
	}

	// Get next available filename
	filePath, err := getNextAvailableFilename(terraformDir, req.Resource, ".tf")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate unique filename",
			"details": err.Error(),
		})
		return
	}

	// Write file with error handling
	if err := os.WriteFile(filePath, []byte(cleanedCode), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to write terraform file",
			"details": err.Error(),
			"path":    filePath,
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, TerraformResponse{
		Message:       "Terraform code generated successfully",
		TerraformCode: cleanedCode,
	})
}

// cleanTerraformCode removes markdown code block markers
func cleanTerraformCode(code string) (string, error) {
	if code == "" {
		return "", fmt.Errorf("input code cannot be empty")
	}

	// Remove ```terraform or ```hcl at the beginning
	re1, err := regexp.Compile(`^` + "```" + `(?:terraform|hcl)\n?`)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex pattern: %w", err)
	}
	cleaned := re1.ReplaceAllString(code, "")
	
	// Remove ``` at the end
	re2, err := regexp.Compile("```$")
	if err != nil {
		return "", fmt.Errorf("failed to compile end regex pattern: %w", err)
	}
	cleaned = re2.ReplaceAllString(cleaned, "")
	
	result := strings.TrimSpace(cleaned)
	if result == "" {
		return "", fmt.Errorf("cleaned code is empty after processing")
	}
	
	return result, nil
}

// getNextAvailableFilename generates a unique filename in the specified directory
func getNextAvailableFilename(dir, resourceText, ext string) (string, error) {
	if dir == "" {
		return "", fmt.Errorf("directory cannot be empty")
	}
	if ext == "" {
		return "", fmt.Errorf("extension cannot be empty")
	}

	// Extract first 5 words from resource and clean them
	words := strings.Fields(resourceText)
	if len(words) > 5 {
		words = words[:5]
	}
	
	baseName := strings.Join(words, "_")
	// Remove non-alphanumeric characters except underscores
	re, err := regexp.Compile(`[^a-zA-Z0-9_]`)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex for filename cleaning: %w", err)
	}
	baseName = re.ReplaceAllString(baseName, "")
	baseName = strings.ToLower(baseName)
	
	if baseName == "" {
		baseName = "generated"
	}

	// Find next available filename (with safety limit to prevent infinite loop)
	maxAttempts := 10000
	for index := 1; index <= maxAttempts; index++ {
		fileName := fmt.Sprintf("%s_%d%s", baseName, index, ext)
		filePath := filepath.Join(dir, fileName)
		
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return filePath, nil
		}
	}
	
	return "", fmt.Errorf("failed to find available filename after %d attempts", maxAttempts)
}
