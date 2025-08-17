package ai

import (
	"context"

	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/services/ai/providers"
)

// Client is the entry point for AI-powered commit message generation
// It follows the same pattern as the custom_commands service
type Client struct {
	c               *helpers.HelperCommon
	contextBuilder  *ContextBuilder
	messageValidator *MessageValidator
}

// NewClient creates a new AI service client
func NewClient(c *helpers.HelperCommon) *Client {
	contextBuilder := NewContextBuilder(c)
	messageValidator := NewMessageValidator()

	return &Client{
		c:               c,
		contextBuilder:  contextBuilder,
		messageValidator: messageValidator,
	}
}

// GenerateCommitMessage generates a commit message based on the current git state
func (c *Client) GenerateCommitMessage(ctx context.Context) (*GenerateResponse, error) {
	// TODO: Implement commit message generation
	// 1. Check if AI is configured and enabled
	// 2. Build context from git state
	// 3. Call provider to generate message
	// 4. Validate and return response
	return nil, nil
}

// GenerateCommitMessageForReword generates an improved commit message for rewording an existing commit
func (c *Client) GenerateCommitMessageForReword(ctx context.Context, existingMessage string, commitSha string) (*GenerateResponse, error) {
	// TODO: Implement reword message generation
	// 1. Get commit diff and context
	// 2. Build prompt with existing message for improvement
	// 3. Call provider to generate improved message
	// 4. Validate and return response
	return nil, nil
}

// IsConfigured returns whether AI is properly configured
func (c *Client) IsConfigured() bool {
	// TODO: Check if AI config is valid and provider is available
	config := c.c.UserConfig().AI
	return config.Enabled && config.APIKey != "" && config.Provider != ""
}

// ValidateConfig validates the current AI configuration
func (c *Client) ValidateConfig() error {
	// TODO: Validate AI configuration
	// 1. Check required fields
	// 2. Validate provider settings
	// 3. Test connectivity if possible
	return nil
}

// getProvider returns the configured AI provider
func (c *Client) getProvider() (Provider, error) {
	// TODO: Initialize and return the appropriate provider based on config
	config := c.c.UserConfig().AI
	
	switch config.Provider {
	case "openai":
		return providers.NewOpenAIProvider(config), nil
	case "github":
		return providers.NewGitHubProvider(config), nil
	default:
		return nil, ErrUnsupportedProvider
	}
}