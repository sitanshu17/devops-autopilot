package models

import "devops-autopilot/utils"

// TerraformRequest represents the request body for terraform generation
type TerraformRequest struct {
	Resource string `json:"resource" binding:"required"`
	Specs    string `json:"specs" binding:"required"`
}

// ValidationRequest represents the request body for terraform validation
type ValidationRequest struct {
	TerraformCode string `json:"terraformCode" binding:"required"`
}

// TerraformResponse represents the response for terraform generation
type TerraformResponse struct {
	Message       string                           `json:"message"`
	TerraformCode string                           `json:"terraformCode"`
	Validation    *utils.TerraformValidationResult `json:"validation,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// ValidationResponse represents the validation-only response
type ValidationResponse struct {
	Message    string                           `json:"message"`
	Validation *utils.TerraformValidationResult `json:"validation"`
}