package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAIConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      AIConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config with all fields",
			config: AIConfig{
				Provider:    "openai",
				BaseURL:     "https://api.openai.com/v1",
				Model:       "gpt-4o-mini",
				APIKeyEnv:   "OPENAI_API_KEY",
				Temperature: 0.2,
				MaxTokens:   300,
				StagedOnly:  true,
				CommitStyle: "conventional",
			},
			expectError: false,
		},
		{
			name: "valid config with minimal fields",
			config: AIConfig{
				Model: "gpt-4o-mini",
			},
			expectError: false,
		},
		{
			name: "missing model",
			config: AIConfig{
				Provider: "openai",
			},
			expectError: true,
			errorMsg:    "ai.model is required",
		},
		{
			name: "empty model",
			config: AIConfig{
				Model: "",
			},
			expectError: true,
			errorMsg:    "ai.model is required",
		},
		{
			name: "temperature too low",
			config: AIConfig{
				Model:       "gpt-4o-mini",
				Temperature: -0.1,
			},
			expectError: true,
			errorMsg:    "ai.temperature must be between 0 and 2",
		},
		{
			name: "temperature too high",
			config: AIConfig{
				Model:       "gpt-4o-mini",
				Temperature: 2.1,
			},
			expectError: true,
			errorMsg:    "ai.temperature must be between 0 and 2",
		},
		{
			name: "valid temperature at boundary",
			config: AIConfig{
				Model:       "gpt-4o-mini",
				Temperature: 0.0,
			},
			expectError: false,
		},
		{
			name: "valid temperature at upper boundary",
			config: AIConfig{
				Model:       "gpt-4o-mini",
				Temperature: 2.0,
			},
			expectError: false,
		},
		{
			name: "negative max tokens",
			config: AIConfig{
				Model:     "gpt-4o-mini",
				MaxTokens: -1,
			},
			expectError: true,
			errorMsg:    "ai.maxTokens must be non-negative",
		},
		{
			name: "max tokens too small",
			config: AIConfig{
				Model:     "gpt-4o-mini",
				MaxTokens: 5,
			},
			expectError: true,
			errorMsg:    "ai.maxTokens must be at least 10 if specified",
		},
		{
			name: "valid max tokens at minimum",
			config: AIConfig{
				Model:     "gpt-4o-mini",
				MaxTokens: 10,
			},
			expectError: false,
		},
		{
			name: "zero max tokens (server decides)",
			config: AIConfig{
				Model:     "gpt-4o-mini",
				MaxTokens: 0,
			},
			expectError: false,
		},
		{
			name: "invalid provider",
			config: AIConfig{
				Model:    "gpt-4o-mini",
				Provider: "invalid-provider",
			},
			expectError: true,
			errorMsg:    "ai.provider must be one of: openai, openrouter, custom",
		},
		{
			name: "valid openai provider",
			config: AIConfig{
				Model:    "gpt-4o-mini",
				Provider: "openai",
			},
			expectError: false,
		},
		{
			name: "valid openrouter provider",
			config: AIConfig{
				Model:    "gpt-4o-mini",
				Provider: "openrouter",
			},
			expectError: false,
		},
		{
			name: "valid custom provider",
			config: AIConfig{
				Model:    "gpt-4o-mini",
				Provider: "custom",
			},
			expectError: false,
		},
		{
			name: "empty provider (defaults to openai)",
			config: AIConfig{
				Model:    "gpt-4o-mini",
				Provider: "",
			},
			expectError: false,
		},
		{
			name: "invalid commit style",
			config: AIConfig{
				Model:       "gpt-4o-mini",
				CommitStyle: "invalid-style",
			},
			expectError: true,
			errorMsg:    "ai.commitStyle must be one of: conventional, plain",
		},
		{
			name: "valid conventional commit style",
			config: AIConfig{
				Model:       "gpt-4o-mini",
				CommitStyle: "conventional",
			},
			expectError: false,
		},
		{
			name: "valid plain commit style",
			config: AIConfig{
				Model:       "gpt-4o-mini",
				CommitStyle: "plain",
			},
			expectError: false,
		},
		{
			name: "empty commit style (defaults to conventional)",
			config: AIConfig{
				Model:       "gpt-4o-mini",
				CommitStyle: "",
			},
			expectError: false,
		},
		{
			name: "multiple validation errors (should return first)",
			config: AIConfig{
				Model:       "", // Missing model
				Temperature: 3.0, // Invalid temperature
				MaxTokens:   -5, // Invalid max tokens
				Provider:    "invalid", // Invalid provider
				CommitStyle: "invalid", // Invalid commit style
			},
			expectError: true,
			errorMsg:    "ai.model is required", // Should return the first error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAIConfigValidateEdgeCases(t *testing.T) {
	// Test with very large values
	t.Run("very large max tokens", func(t *testing.T) {
		config := AIConfig{
			Model:     "gpt-4o-mini",
			MaxTokens: 1000000,
		}
		err := config.Validate()
		assert.NoError(t, err)
	})

	// Test with very precise temperature values
	t.Run("precise temperature values", func(t *testing.T) {
		config := AIConfig{
			Model:       "gpt-4o-mini",
			Temperature: 1.9999999,
		}
		err := config.Validate()
		assert.NoError(t, err)
	})

	// Test with whitespace in string fields
	t.Run("whitespace in model name", func(t *testing.T) {
		config := AIConfig{
			Model: "  gpt-4o-mini  ",
		}
		// Note: The validation doesn't trim whitespace, so this should pass
		// In a real implementation, you might want to add trimming
		err := config.Validate()
		assert.NoError(t, err)
	})
}

func TestAIConfigDefaults(t *testing.T) {
	// Test that the default configuration from GetDefaultConfig is valid
	defaultConfig := GetDefaultConfig()
	err := defaultConfig.AI.Validate()
	
	// The default config has an empty model, so it should fail validation
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ai.model is required")
	
	// But if we set a model, it should be valid
	defaultConfig.AI.Model = "gpt-4o-mini"
	err = defaultConfig.AI.Validate()
	assert.NoError(t, err)
}
