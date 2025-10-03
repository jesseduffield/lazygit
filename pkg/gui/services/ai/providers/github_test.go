package providers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
)

func TestGitHubProvider_Name(t *testing.T) {
	provider := NewGitHubProvider(config.AIConfig{})
	if provider.Name() != "github" {
		t.Errorf("Expected provider name to be 'github', got '%s'", provider.Name())
	}
}

func TestGitHubProvider_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  config.AIConfig
		wantErr bool
	}{
		{
			name:    "valid config without API key (OAuth used)",
			config:  config.AIConfig{},
			wantErr: false,
		},
		{
			name: "valid config with all fields",
			config: config.AIConfig{
				Provider:    "github",
				Model:       "gpt-4",
				Temperature: 0.7,
				MaxTokens:   500,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewGitHubProvider(tt.config)
			err := provider.ValidateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GitHubProvider.ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitHubProvider_buildRequest(t *testing.T) {
	config := config.AIConfig{
		Model:       "gpt-4",
		Temperature: 0.8,
		MaxTokens:   1000,
	}
	provider := NewGitHubProvider(config)

	prompt := "Add user authentication functionality"
	request := provider.buildRequest(prompt)

	if request.Model != "gpt-4" {
		t.Errorf("Expected model to be 'gpt-4', got '%s'", request.Model)
	}

	if request.Temperature != 0.8 {
		t.Errorf("Expected temperature to be 0.8, got %f", request.Temperature)
	}

	if request.MaxTokens != 1000 {
		t.Errorf("Expected max tokens to be 1000, got %d", request.MaxTokens)
	}

	if len(request.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(request.Messages))
	}

	if request.Messages[0].Role != "system" {
		t.Errorf("Expected first message role to be 'system', got '%s'", request.Messages[0].Role)
	}

	if request.Messages[1].Role != "user" {
		t.Errorf("Expected second message role to be 'user', got '%s'", request.Messages[1].Role)
	}

	if request.Messages[1].Content != "Generate a git commit message for the following changes:\n\nAdd user authentication functionality" {
		t.Errorf("Unexpected user message content: %s", request.Messages[1].Content)
	}
}

func TestGitHubProvider_buildRequest_defaults(t *testing.T) {
	config := config.AIConfig{} // Empty config to test defaults
	provider := NewGitHubProvider(config)

	request := provider.buildRequest("test prompt")

	if request.Model != "gpt-4" {
		t.Errorf("Expected default model to be 'gpt-4', got '%s'", request.Model)
	}

	if request.Temperature != 0.7 {
		t.Errorf("Expected default temperature to be 0.7, got %f", request.Temperature)
	}

	if request.MaxTokens != 500 {
		t.Errorf("Expected default max tokens to be 500, got %d", request.MaxTokens)
	}

	if request.TopP != 1.0 {
		t.Errorf("Expected TopP to be 1.0, got %f", request.TopP)
	}

	if request.N != 1 {
		t.Errorf("Expected N to be 1, got %d", request.N)
	}

	if request.Stream != false {
		t.Errorf("Expected Stream to be false, got %t", request.Stream)
	}
}

func TestGitHubProvider_extractMessage(t *testing.T) {
	tests := []struct {
		name     string
		response *copilotResponse
		want     string
		wantErr  bool
	}{
		{
			name: "valid response",
			response: &copilotResponse{
				Choices: []copilotChoice{
					{
						Message: copilotMessage{
							Content: "feat: add user authentication",
						},
					},
				},
			},
			want:    "feat: add user authentication",
			wantErr: false,
		},
		{
			name: "response with quotes and whitespace",
			response: &copilotResponse{
				Choices: []copilotChoice{
					{
						Message: copilotMessage{
							Content: "  \"fix: resolve login bug\"  \n",
						},
					},
				},
			},
			want:    "fix: resolve login bug",
			wantErr: false,
		},
		{
			name: "empty choices",
			response: &copilotResponse{
				Choices: []copilotChoice{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty message content",
			response: &copilotResponse{
				Choices: []copilotChoice{
					{
						Message: copilotMessage{
							Content: "",
						},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
	}

	provider := NewGitHubProvider(config.AIConfig{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := provider.extractMessage(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitHubProvider.extractMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GitHubProvider.extractMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
