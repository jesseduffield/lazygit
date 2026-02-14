package helpers

import (
	"context"
	"errors"
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/ai"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type AIHelper struct {
	c         *HelperCommon
	aiManager *ai.AIManager
	iGetHelpers
}

type iGetHelpers interface {
	Helpers() *Helpers
}

func NewAIHelper(c *HelperCommon, aiManager *ai.AIManager, iGetHelpers iGetHelpers) *AIHelper {
	return &AIHelper{
		c:           c,
		aiManager:   aiManager,
		iGetHelpers: iGetHelpers,
	}
}

func (self *AIHelper) GenerateCommitMessage() error {
	return self.c.WithWaitingStatus(self.c.Tr.AIGeneratingCommitMessage, func(gocui.Task) error {
		ctx := context.Background()

		self.c.LogAction(self.c.Tr.AIGenerateCommitMessageAction)

		result, err := self.aiManager.GenerateCommitMessage(ctx)
		if err != nil {
			if errors.Is(err, ai.ErrNotConfigured) {
				return fmt.Errorf(self.c.Tr.AINotConfigured)
			}
			if errors.Is(err, ai.ErrProviderNotSupported) {
				aiConfig := self.c.UserConfig().AI
				provider := ""
				if aiConfig != nil {
					provider = aiConfig.Provider
				}
				message := utils.ResolvePlaceholderString(
					self.c.Tr.AIProviderNotSupported,
					map[string]string{
						"provider": provider,
					},
				)
				return fmt.Errorf(message)
			}
			if errors.Is(err, ai.ErrNoStagedChanges) {
				return fmt.Errorf(self.c.Tr.AINoStagedChanges)
			}
			return fmt.Errorf("%s: %w", self.c.Tr.AIGenerationFailed, err)
		}

		self.c.LogCommand("Prompt:\n"+result.Prompt, false)
		self.c.LogCommand("Response:\n"+result.RawResponse, false)

		message := result.Message
		if result.Description != "" {
			message = result.Message + "\n\n" + result.Description
		}

		self.Helpers().Commits.SetMessageAndDescriptionInView(message)

		return nil
	})
}
