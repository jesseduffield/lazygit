package controllers

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type ShellCommandAction struct {
	c *ControllerCommon
}

func (self *ShellCommandAction) Call() error {
	self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.ShellCommand,
		FindSuggestionsFunc: self.GetShellCommandsHistorySuggestionsFunc(),
		AllowEditSuggestion: true,
		PreserveWhitespace:  true,
		HandleConfirm: func(command string) error {
			if self.shouldSaveCommand(command) {
				self.c.UpdateAppState(func(state *config.AppState) {
					state.ShellCommandsHistory = utils.Limit(
						lo.Uniq(append([]string{command}, state.ShellCommandsHistory...)),
						1000,
					)
				})
			}

			self.c.LogAction(self.c.Tr.Actions.CustomCommand)
			return self.c.RunSubprocessAndRefresh(
				self.c.OS().Cmd.NewShell(command, self.c.UserConfig().OS.ShellFunctionsFile),
			)
		},
		HandleDeleteSuggestion: func(index int) error {
			// index is the index in the _filtered_ list of suggestions, not in
			// the full list. We keep the items in the ShellCommandsHistory
			// unique, so we can remove the item by its value.
			item := self.c.Contexts().Suggestions.GetItems()[index].Value
			self.c.UpdateAppState(func(state *config.AppState) {
				state.ShellCommandsHistory = lo.Without(state.ShellCommandsHistory, item)
			})
			self.c.Contexts().Suggestions.RefreshSuggestions()
			return nil
		},
	})

	return nil
}

func (self *ShellCommandAction) GetShellCommandsHistorySuggestionsFunc() func(string) []*types.Suggestion {
	return func(input string) []*types.Suggestion {
		history := self.c.ReadAppState().ShellCommandsHistory

		return helpers.FilterFunc(history, self.c.UserConfig().Gui.UseFuzzySearch())(input)
	}
}

// this mimics the shell functionality `ignorespace`
// which doesn't save a command to history if it starts with a space
func (self *ShellCommandAction) shouldSaveCommand(command string) bool {
	return !strings.HasPrefix(command, " ")
}
