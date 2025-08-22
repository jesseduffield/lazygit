package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
)

// OpenAIProvider implements the Provider interface for OpenAI API
type OpenAIProvider struct {
	config     config.AIConfig
	httpClient *http.Client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config config.AIConfig) *OpenAIProvider {
	return &OpenAIProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateMessage generates a commit message using OpenAI API
func (p *OpenAIProvider) GenerateMessage(prompt string) (string, error) {
	// TODO: Implement OpenAI API call
	// 1. Build request payload
	// 2. Make HTTP request to OpenAI API
	// 3. Parse response
	// 4. Return generated message

	return "", fmt.Errorf("OpenAI provider not implemented yet")
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// ValidateConfig validates the OpenAI configuration
func (p *OpenAIProvider) ValidateConfig() error {
	// TODO: Validate OpenAI configuration
	// 1. Check API key is present
	// 2. Validate model name
	// 3. Check base URL format
	// 4. Optionally test API connectivity

	if p.config.APIKey == "" {
		return fmt.Errorf("OpenAI API key is required")
	}

	if p.config.Model == "" {
		return fmt.Errorf("OpenAI model is required")
	}

	return nil
}

// OpenAI API request/response structures
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	Stream      bool            `json:"stream"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []openAIChoice `json:"choices"`
	Usage   openAIUsage    `json:"usage"`
	Error   *openAIError   `json:"error,omitempty"`
}

type openAIChoice struct {
	Index        int           `json:"index"`
	Message      openAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

type openAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type openAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// buildRequest creates the OpenAI API request payload
func (p *OpenAIProvider) buildRequest(prompt string) *openAIRequest {
	// TODO: Build proper request with system and user messages
	return &openAIRequest{
		Model: p.config.Model,
		Messages: []openAIMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant that generates concise, clear git commit messages based on code changes.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   p.config.MaxTokens,
		Temperature: p.config.Temperature,
		Stream:      false,
	}
}

// makeAPICall makes the HTTP request to OpenAI API
func (p *OpenAIProvider) makeAPICall(ctx context.Context, request *openAIRequest) (*openAIResponse, error) {
	// TODO: Implement actual API call
	// 1. Serialize request to JSON
	// 2. Create HTTP request with proper headers
	// 3. Make request with context
	// 4. Parse response
	// 5. Handle errors

	// Serialize request
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := p.config.BaseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	// Make request (stub for now)
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var apiResponse openAIResponse
	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if apiResponse.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", apiResponse.Error.Message)
	}

	return &apiResponse, nil
}

// extractMessage extracts the generated message from OpenAI response
func (p *OpenAIProvider) extractMessage(response *openAIResponse) (string, error) {
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	choice := response.Choices[0]
	message := choice.Message.Content

	if message == "" {
		return "", fmt.Errorf("empty message in response")
	}

	return message, nil
}
