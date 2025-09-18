package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientFromConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      config.AIConfig
		envVars     map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config with OpenAI",
			config: config.AIConfig{
				Provider:    "openai",
				Model:       "gpt-4o-mini",
				Temperature: 0.2,
				MaxTokens:   300,
			},
			envVars: map[string]string{
				"OPENAI_API_KEY": "test-key",
			},
			expectError: false,
		},
		{
			name: "valid config with custom provider",
			config: config.AIConfig{
				Provider:    "custom",
				BaseURL:     "https://api.example.com/v1",
				Model:       "custom-model",
				APIKeyEnv:   "CUSTOM_API_KEY",
				Temperature: 0.5,
				MaxTokens:   500,
			},
			envVars: map[string]string{
				"CUSTOM_API_KEY": "custom-key",
			},
			expectError: false,
		},
		{
			name: "missing model",
			config: config.AIConfig{
				Provider: "openai",
			},
			envVars: map[string]string{
				"OPENAI_API_KEY": "test-key",
			},
			expectError: true,
			errorMsg:    "ai.model is required",
		},
		{
			name: "invalid temperature",
			config: config.AIConfig{
				Provider:    "openai",
				Model:       "gpt-4o-mini",
				Temperature: 3.0,
			},
			envVars: map[string]string{
				"OPENAI_API_KEY": "test-key",
			},
			expectError: true,
			errorMsg:    "ai.temperature must be between 0 and 2",
		},
		{
			name: "missing API key",
			config: config.AIConfig{
				Provider: "openai",
				Model:    "gpt-4o-mini",
			},
			envVars:     map[string]string{},
			expectError: true,
			errorMsg:    "missing API key; set OPENAI_API_KEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			client, err := NewClientFromConfig(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tt.config.Model, client.model)
				assert.Equal(t, tt.config.Temperature, client.temp)
				assert.Equal(t, tt.config.MaxTokens, client.maxTokens)
			}
		})
	}
}

func TestDefaultBaseURL(t *testing.T) {
	tests := []struct {
		provider string
		expected string
	}{
		{"openai", "https://api.openai.com/v1"},
		{"openrouter", "https://openrouter.ai/api/v1"},
		{"", "https://api.openai.com/v1"},
		{"custom-url", "custom-url"},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			result := defaultBaseURL(tt.provider)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultAPIKeyEnv(t *testing.T) {
	tests := []struct {
		provider string
		expected string
	}{
		{"openrouter", "OPENROUTER_API_KEY"},
		{"openai", "OPENAI_API_KEY"},
		{"", "OPENAI_API_KEY"},
		{"custom", "OPENAI_API_KEY"},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			result := defaultAPIKeyEnv(tt.provider)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateCommitMessage(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/chat/completions", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))

		// Decode request to verify it's properly formatted
		var req ChatRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "test-model", req.Model)
		assert.Equal(t, 0.2, req.Temperature)
		assert.Equal(t, 300, req.MaxTokens)
		assert.Len(t, req.Messages, 2)
		assert.Equal(t, "system", req.Messages[0].Role)
		assert.Equal(t, "user", req.Messages[1].Role)

		// Mock response
		response := ChatResponse{
			ID:      "test-id",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Choices: []ChatChoice{
				{
					Index:        0,
					FinishReason: "stop",
					Message: ChatMessage{
						Role:    "assistant",
						Content: "feat: add new feature\n\nThis commit adds a new feature to improve user experience.",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		http:      &http.Client{Timeout: 10 * time.Second},
		baseURL:   server.URL,
		apiKey:    "test-key",
		model:     "test-model",
		temp:      0.2,
		maxTokens: 300,
	}

	ctx := context.Background()
	diff := "diff --git a/file.go b/file.go\n+func newFeature() {}"
	
	subject, body, err := client.GenerateCommitMessage(ctx, diff, "conventional", 72)
	
	assert.NoError(t, err)
	assert.Equal(t, "feat: add new feature", subject)
	assert.Equal(t, "This commit adds a new feature to improve user experience.", body)
}

func TestGenerateCommitMessageWithRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			// Fail first two attempts
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		
		// Succeed on third attempt
		response := ChatResponse{
			Choices: []ChatChoice{
				{
					Message: ChatMessage{
						Role:    "assistant",
						Content: "fix: resolve issue",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		http:      &http.Client{Timeout: 10 * time.Second},
		baseURL:   server.URL,
		apiKey:    "test-key",
		model:     "test-model",
		temp:      0.2,
		maxTokens: 300,
	}

	ctx := context.Background()
	diff := "diff --git a/file.go b/file.go\n-bug\n+fix"
	
	subject, body, err := client.GenerateCommitMessage(ctx, diff, "conventional", 72)
	
	assert.NoError(t, err)
	assert.Equal(t, "fix: resolve issue", subject)
	assert.Equal(t, "", body)
	assert.Equal(t, 3, attempts) // Should have retried twice
}

func TestGenerateCommitMessageLongSubject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ChatResponse{
			Choices: []ChatChoice{
				{
					Message: ChatMessage{
						Role:    "assistant",
						Content: strings.Repeat("a", 100), // Very long subject
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		http:      &http.Client{Timeout: 10 * time.Second},
		baseURL:   server.URL,
		apiKey:    "test-key",
		model:     "test-model",
		temp:      0.2,
		maxTokens: 300,
	}

	ctx := context.Background()
	subject, _, err := client.GenerateCommitMessage(ctx, "diff", "plain", 72)
	
	assert.NoError(t, err)
	assert.Len(t, []rune(subject), 72) // Should be truncated to 72 characters
}
