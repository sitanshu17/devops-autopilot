package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TerraformValidationResult holds the result of terraform validation
type TerraformValidationResult struct {
	IsValid   bool     `json:"isValid"`
	Errors    []string `json:"errors,omitempty"`
	Warnings  []string `json:"warnings,omitempty"`
	Output    string   `json:"output,omitempty"`
	ExecTime  int64    `json:"execTime"` // milliseconds
}

// ValidateTerraformCode validates terraform code using local terraform CLI
func ValidateTerraformCode(terraformCode string) (*TerraformValidationResult, error) {
	startTime := time.Now()
	
	// Check if terraform CLI is available
	if !isTerraformInstalled() {
		return &TerraformValidationResult{
			IsValid: false,
			Errors:  []string{"Terraform CLI is not installed or not available in PATH"},
			ExecTime: time.Since(startTime).Milliseconds(),
		}, nil
	}

	// Create temporary directory for validation
	tempDir, err := createTempTerraformDir(terraformCode)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer cleanupTempDir(tempDir)

	// Run terraform init (required before validate)
	initResult, err := runTerraformInit(tempDir)
	if err != nil {
		return &TerraformValidationResult{
			IsValid: false,
			Errors:  []string{fmt.Sprintf("Terraform init failed: %s", err.Error())},
			Output:  initResult,
			ExecTime: time.Since(startTime).Milliseconds(),
		}, nil
	}

	// Run terraform validate
	validateResult, err := runTerraformValidate(tempDir)
	execTime := time.Since(startTime).Milliseconds()

	if err != nil {
		// Parse terraform validation errors from the actual output
		errors := parseTerraformErrors(validateResult)
		return &TerraformValidationResult{
			IsValid: false,
			Errors:  errors,
			Output:  validateResult,
			ExecTime: execTime,
		}, nil
	}

	return &TerraformValidationResult{
		IsValid:  true,
		Output:   validateResult,
		ExecTime: execTime,
	}, nil
}

// isTerraformInstalled checks if terraform CLI is available
func isTerraformInstalled() bool {
	_, err := exec.LookPath("terraform")
	return err == nil
}

// createTempTerraformDir creates a temporary directory with the terraform code
func createTempTerraformDir(terraformCode string) (string, error) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "terraform_validate_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Write terraform code to main.tf
	mainTfPath := filepath.Join(tempDir, "main.tf")
	err = ioutil.WriteFile(mainTfPath, []byte(terraformCode), 0644)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to write terraform file: %w", err)
	}

	return tempDir, nil
}

// runTerraformInit runs terraform init in the given directory
func runTerraformInit(dir string) (string, error) {
	cmd := exec.Command("terraform", "init", "-no-color")
	cmd.Dir = dir
	
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	
	if err != nil {
		return outputStr, fmt.Errorf("terraform init failed: %w", err)
	}
	
	return outputStr, nil
}

// runTerraformValidate runs terraform validate in the given directory
func runTerraformValidate(dir string) (string, error) {
	cmd := exec.Command("terraform", "validate", "-no-color", "-json")
	cmd.Dir = dir
	
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	
	if err != nil {
		return outputStr, fmt.Errorf("terraform validate failed: %w", err)
	}
	
	return outputStr, nil
}

// parseTerraformErrors extracts error messages from terraform output
func parseTerraformErrors(output string) []string {
	var errors []string
	
	// Try to parse JSON output first (terraform validate -json)
	if strings.Contains(output, `"valid": false`) {
		errors = parseJSONErrors(output)
		if len(errors) > 0 {
			return errors
		}
	}
	
	// Fallback to line-by-line parsing
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "terraform validate failed:") {
			errors = append(errors, line)
		}
	}
	
	if len(errors) == 0 {
		errors = []string{output}
	}
	
	return errors
}

// TerraformDiagnostic represents a single diagnostic from terraform validate
type TerraformDiagnostic struct {
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
	Detail   string `json:"detail"`
	Range    struct {
		Filename string `json:"filename"`
		Start    struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"start"`
	} `json:"range"`
}

// TerraformValidateOutput represents the JSON output from terraform validate
type TerraformValidateOutput struct {
	Valid       bool                  `json:"valid"`
	ErrorCount  int                   `json:"error_count"`
	Diagnostics []TerraformDiagnostic `json:"diagnostics"`
}

// parseJSONErrors extracts errors from terraform's JSON output
func parseJSONErrors(jsonOutput string) []string {
	var errors []string
	var validateOutput TerraformValidateOutput
	
	// Try to parse the JSON
	if err := json.Unmarshal([]byte(jsonOutput), &validateOutput); err != nil {
		// If JSON parsing fails, fall back to simple string extraction
		return []string{fmt.Sprintf("Failed to parse validation output: %s", jsonOutput)}
	}
	
	// Extract meaningful error messages
	for _, diagnostic := range validateOutput.Diagnostics {
		if diagnostic.Severity == "error" {
			errorMsg := fmt.Sprintf("Line %d: %s - %s", 
				diagnostic.Range.Start.Line, 
				diagnostic.Summary, 
				diagnostic.Detail)
			errors = append(errors, errorMsg)
		}
	}
	
	// If no errors found in diagnostics but validation failed, show generic error
	if len(errors) == 0 && !validateOutput.Valid {
		errors = []string{fmt.Sprintf("Validation failed with %d errors", validateOutput.ErrorCount)}
	}
	
	return errors
}

// cleanupTempDir removes the temporary directory
func cleanupTempDir(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to cleanup temp dir %s: %v\n", dir, err)
	}
}
