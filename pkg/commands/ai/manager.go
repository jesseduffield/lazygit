package ai

import (
	"context"
	"errors"
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/config"
)

var (
	ErrNotConfigured        = errors.New("AI not configured")
	ErrProviderNotSupported = errors.New("provider not supported")
	ErrNoStagedChanges      = errors.New("no staged changes")
)

type AIManager struct {
	gitCmd         *commands.GitCommand
	config         *config.AppConfig
	provider       Provider
	contextBuilder *ContextBuilder
}

func NewAIManager(gitCmd *commands.GitCommand, config *config.AppConfig) *AIManager {
	manager := &AIManager{
		gitCmd:         gitCmd,
		config:         config,
		contextBuilder: NewContextBuilder(),
	}

	manager.initProvider()
	return manager
}

func (m *AIManager) initProvider() {
	aiConfig := m.config.GetUserConfig().AI
	if aiConfig == nil {
		return
	}

	switch ProviderType(aiConfig.Provider) {
	case ProviderOllama:
		m.provider = NewOllamaProvider(aiConfig.BaseURL, aiConfig.Model)
	}
}

func (m *AIManager) IsEnabled() bool {
	return m.config.GetUserConfig().AI != nil && m.provider != nil
}

func (m *AIManager) GenerateCommitMessage(ctx context.Context) (*GenerateResponse, error) {
	aiConfig := m.config.GetUserConfig().AI
	if aiConfig == nil {
		return nil, ErrNotConfigured
	}

	if m.provider == nil {
		return nil, ErrProviderNotSupported
	}

	diff, err := m.gitCmd.Diff.GetDiff(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get diff: %w", err)
	}

	if diff == "" {
		return nil, ErrNoStagedChanges
	}

	branchName, err := m.gitCmd.Branch.CurrentBranchName()
	if err != nil {
		branchName = "unknown"
	}

	filteredDiff := m.contextBuilder.BuildContext(diff)

	req := GenerateRequest{
		IncludeDescription:     aiConfig.GenerateDescription,
		UseConventionalCommits: aiConfig.UseConventionalCommits,
		Diff:                   filteredDiff,
		BranchName:             branchName,
	}

	return m.provider.Generate(ctx, req)
}
