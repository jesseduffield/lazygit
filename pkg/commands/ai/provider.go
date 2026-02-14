package ai

import (
	"context"
)

type ProviderType string

const (
	ProviderOllama ProviderType = "ollama"
	ProviderOpenAI ProviderType = "openai"
	ProviderGemini ProviderType = "gemini"
	ProviderClaude ProviderType = "claude"
)

type GenerateRequest struct {
	Prompt                 string
	IncludeDescription     bool
	UseConventionalCommits bool
	Diff                   string
	BranchName             string
}

type GenerateResponse struct {
	Message     string
	Description string
	Error       error
	Prompt      string
	RawResponse string
}

type Provider interface {
	Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)
	Name() ProviderType
}
