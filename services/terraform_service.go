package services

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"devops-autopilot/utils"
)

// TerraformService handles terraform-related business logic
type TerraformService struct{}

// NewTerraformService creates a new terraform service
func NewTerraformService() *TerraformService {
	return &TerraformService{}
}

// GenerateAndValidate generates terraform code and validates it
func (s *TerraformService) GenerateAndValidate(resource, specs string) (string, *utils.TerraformValidationResult, error) {
	// Generate terraform code using OpenAI
	tfCode, err := utils.GenerateTerraformCode(resource, specs)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate terraform code: %w", err)
	}

	// Validate generated code is not empty
	if strings.TrimSpace(tfCode) == "" {
		return "", nil, fmt.Errorf("generated terraform code is empty")
	}

	// Clean the code (remove markdown code block markers)
	cleanedCode, err := s.CleanTerraformCode(tfCode)
	if err != nil {
		return "", nil, fmt.Errorf("failed to clean terraform code: %w", err)
	}

	// Validate the generated Terraform code
	validation, err := utils.ValidateTerraformCode(cleanedCode)
	if err != nil {
		return cleanedCode, nil, fmt.Errorf("failed to validate terraform code: %w", err)
	}

	return cleanedCode, validation, nil
}

// GenerateAndValidateWithCopilot generates terraform code using GitHub Copilot and validates it
func (s *TerraformService) GenerateAndValidateWithCopilot(resource, specs string) (string, *utils.TerraformValidationResult, error) {
	// Generate terraform code using GitHub Copilot
	tfCode, err := utils.GenerateTerraformCodeWithCopilot(resource, specs)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate terraform code with GitHub Copilot: %w", err)
	}

	// Validate generated code is not empty
	if strings.TrimSpace(tfCode) == "" {
		return "", nil, fmt.Errorf("generated terraform code is empty")
	}

	// Clean the code (remove markdown code block markers)
	cleanedCode, err := s.CleanTerraformCode(tfCode)
	if err != nil {
		return "", nil, fmt.Errorf("failed to clean terraform code: %w", err)
	}

	// Validate the generated Terraform code
	validation, err := utils.ValidateTerraformCode(cleanedCode)
	if err != nil {
		return cleanedCode, nil, fmt.Errorf("failed to validate terraform code: %w", err)
	}

	return cleanedCode, validation, nil
}

// SaveTerraformFile saves terraform code to a file with provider prefix
func (s *TerraformService) SaveTerraformFile(code, resource, provider string) (string, error) {
	// Ensure tf-generated-files directory exists
	terraformDir := "tf-generated-files"
	if err := os.MkdirAll(terraformDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create tf-generated-files directory: %w", err)
	}

	// Get next available filename with provider prefix
	filePath, err := s.GetNextAvailableFilename(terraformDir, resource, provider, ".tf")
	if err != nil {
		return "", fmt.Errorf("failed to generate unique filename: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return "", fmt.Errorf("failed to write terraform file: %w", err)
	}

	return filePath, nil
}

// CleanTerraformCode removes markdown code block markers
func (s *TerraformService) CleanTerraformCode(code string) (string, error) {
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

// GetNextAvailableFilename generates a unique filename with provider prefix in the specified directory
func (s *TerraformService) GetNextAvailableFilename(dir, resourceText, provider, ext string) (string, error) {
	if dir == "" {
		return "", fmt.Errorf("directory cannot be empty")
	}
	if ext == "" {
		return "", fmt.Errorf("extension cannot be empty")
	}
	if provider == "" {
		return "", fmt.Errorf("provider cannot be empty")
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

	// Add provider prefix to the base name
	baseNameWithProvider := fmt.Sprintf("%s_%s", provider, baseName)

	// Find next available filename (with safety limit to prevent infinite loop)
	maxAttempts := 10000
	for index := 1; index <= maxAttempts; index++ {
		fileName := fmt.Sprintf("%s_%d%s", baseNameWithProvider, index, ext)
		filePath := filepath.Join(dir, fileName)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return filePath, nil
		}
	}

	return "", fmt.Errorf("failed to find available filename after %d attempts", maxAttempts)
}
