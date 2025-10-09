package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

var openaiClient *openai.Client

// InitOpenAI initializes the OpenAI client
func InitOpenAI() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}

	openaiClient = openai.NewClient(apiKey)
	log.Println("OpenAI client initialized successfully")
}

// GenerateTerraformCode generates Terraform code using OpenAI API
func GenerateTerraformCode(resource, specs string) (string, error) {
	// Validate inputs
	if openaiClient == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}

	if strings.TrimSpace(resource) == "" {
		return "", fmt.Errorf("resource cannot be empty")
	}

	if strings.TrimSpace(specs) == "" {
		return "", fmt.Errorf("specs cannot be empty")
	}

	log.Printf("Generating Terraform code for resource: %s with specs: %s", resource, specs)

	prompt := fmt.Sprintf(`
You are a Terraform expert. Generate Terraform code to provision the following:

Resource: %s
Specs: %s

Only output valid Terraform code inside one block. Do not explain anything.
`, resource, specs)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.2,
		MaxTokens:   2000, // Limit response size
	}

	// Create context with timeout for the API call
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Printf("Error calling OpenAI API: %v", err)
		return "", fmt.Errorf("failed to generate terraform code: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices from OpenAI API")
	}

	content := resp.Choices[0].Message.Content
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("OpenAI returned empty content")
	}

	log.Printf("Successfully generated Terraform code (%d characters)", len(content))
	return content, nil
}
