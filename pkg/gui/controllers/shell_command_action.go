package controllers

import (
	"slices"
	"strings"

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
		HandleConfirm: func(command string) error {
			if self.shouldSaveCommand(command) {
				self.c.GetAppState().ShellCommandsHistory = utils.Limit(
					lo.Uniq(append([]string{command}, self.c.GetAppState().ShellCommandsHistory...)),
					1000,
				)
			}

			self.c.SaveAppStateAndLogError()

			self.c.LogAction(self.c.Tr.Actions.CustomCommand)
			return self.c.RunSubprocessAndRefresh(
				self.c.OS().Cmd.NewInteractiveShell(command),
			)
		},
		HandleDeleteSuggestion: func(index int) error {
			// index is the index in the _filtered_ list of suggestions, so we
			// need to map it back to the full list. There's no really good way
			// to do this, but fortunately we keep the items in the
			// ShellCommandsHistory unique, which allows us to simply search
			// for it by string.
			item := self.c.Contexts().Suggestions.GetItems()[index].Value
			fullIndex := lo.IndexOf(self.c.GetAppState().ShellCommandsHistory, item)
			if fullIndex == -1 {
				// Should never happen, but better be safe
				return nil
			}

			self.c.GetAppState().ShellCommandsHistory = slices.Delete(
				self.c.GetAppState().ShellCommandsHistory, fullIndex, fullIndex+1)
			self.c.SaveAppStateAndLogError()
			self.c.Contexts().Suggestions.RefreshSuggestions()
			return nil
		},
	})

	return nil
}

func (self *ShellCommandAction) GetShellCommandsHistorySuggestionsFunc() func(string) []*types.Suggestion {
	return func(input string) []*types.Suggestion {
		history := self.c.GetAppState().ShellCommandsHistory

		return helpers.FilterFunc(history, self.c.UserConfig().Gui.UseFuzzySearch())(input)
	}
}

// this mimics the shell functionality `ignorespace`
// which doesn't save a command to history if it starts with a space
func (self *ShellCommandAction) shouldSaveCommand(command string) bool {
	return !strings.HasPrefix(command, " ")
}
