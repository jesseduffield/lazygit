package controllers

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type CustomCommandAction struct {
	c *ControllerCommon
}

func (self *CustomCommandAction) Call() error {
	return self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.CustomCommand,
		FindSuggestionsFunc: self.GetCustomCommandsHistorySuggestionsFunc(),
		HandleConfirm: func(command string) error {
			if self.shouldSaveCommand(command) {
				self.c.GetAppState().CustomCommandsHistory = utils.Limit(
					lo.Uniq(append(self.c.GetAppState().CustomCommandsHistory, command)),
					1000,
				)
			}

			self.c.SaveAppStateAndLogError()

			self.c.LogAction(self.c.Tr.Actions.CustomCommand)
			return self.c.RunSubprocessAndRefresh(
				self.c.OS().Cmd.NewShell(command),
			)
		},
	})
}

func (self *CustomCommandAction) GetCustomCommandsHistorySuggestionsFunc() func(string) []*types.Suggestion {
	// reversing so that we display the latest command first
	history := lo.Reverse(self.c.GetAppState().CustomCommandsHistory)

	return helpers.FuzzySearchFunc(history)
}

// this mimics the shell functionality `ignorespace`
// which doesn't save a command to history if it starts with a space
func (self *CustomCommandAction) shouldSaveCommand(command string) bool {
	return !strings.HasPrefix(command, " ")
}
