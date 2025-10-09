package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// GitHubMessage represents a message in the chat completion request
type GitHubMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GitHubChatRequest represents the request to GitHub Models API
type GitHubChatRequest struct {
	Messages    []GitHubMessage `json:"messages"`
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature"`
}

// GitHubChoice represents a choice in the response
type GitHubChoice struct {
	Message GitHubMessage `json:"message"`
}

// GitHubChatResponse represents the response from GitHub Models API
type GitHubChatResponse struct {
	Choices []GitHubChoice `json:"choices"`
}

var githubClient *http.Client

// InitGitHub initializes the GitHub client
func InitGitHub() {
	githubClient = &http.Client{
		Timeout: 30 * time.Second,
	}

	// Validate GitHub token exists
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Println("Warning: GITHUB_TOKEN environment variable is not set - GitHub Copilot API will not work")
		return
	}

	log.Println("GitHub client initialized successfully")
}

// GenerateTerraformCodeWithCopilot generates Terraform code using GitHub Models API
func GenerateTerraformCodeWithCopilot(resource, specs string) (string, error) {
	// Validate inputs
	if githubClient == nil {
		return "", fmt.Errorf("GitHub client not initialized")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	if strings.TrimSpace(resource) == "" {
		return "", fmt.Errorf("resource cannot be empty")
	}

	if strings.TrimSpace(specs) == "" {
		return "", fmt.Errorf("specs cannot be empty")
	}

	log.Printf("Generating Terraform code using GitHub Copilot for resource: %s with specs: %s", resource, specs)

	prompt := fmt.Sprintf(`You are a Terraform expert. Generate Terraform code to provision the following:

Resource: %s
Specs: %s

Only output valid Terraform code inside one block. Do not explain anything.
The code should be production-ready and follow best practices.`, resource, specs)

	// Prepare the request
	request := GitHubChatRequest{
		Messages: []GitHubMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:       "gpt-4o-mini", // Using GPT-4o-mini for cost efficiency
		MaxTokens:   2000,
		Temperature: 0.2,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "https://models.inference.ai.azure.com/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Make the request
	resp, err := githubClient.Do(req)
	if err != nil {
		log.Printf("Error calling GitHub Models API: %v", err)
		return "", fmt.Errorf("failed to call GitHub Models API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		log.Printf("GitHub Models API returned status %d: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("GitHub Models API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response GitHubChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices from GitHub Models API")
	}

	content := response.Choices[0].Message.Content
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("GitHub Models API returned empty content")
	}

	log.Printf("Successfully generated Terraform code using GitHub Copilot (%d characters)", len(content))
	return content, nil
}
