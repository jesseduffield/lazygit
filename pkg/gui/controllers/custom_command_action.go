package controllers

import (
	"slices"
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
		AllowEditSuggestion: true,
		HandleConfirm: func(command string) error {
			if self.shouldSaveCommand(command) {
				self.c.GetAppState().CustomCommandsHistory = utils.Limit(
					lo.Uniq(append([]string{command}, self.c.GetAppState().CustomCommandsHistory...)),
					1000,
				)
			}

			self.c.SaveAppStateAndLogError()

			self.c.LogAction(self.c.Tr.Actions.CustomCommand)
			return self.c.RunSubprocessAndRefresh(
				self.c.OS().Cmd.NewShell(command),
			)
		},
		HandleDeleteSuggestion: func(index int) error {
			// index is the index in the _filtered_ list of suggestions, so we
			// need to map it back to the full list. There's no really good way
			// to do this, but fortunately we keep the items in the
			// CustomCommandsHistory unique, which allows us to simply search
			// for it by string.
			item := self.c.Contexts().Suggestions.GetItems()[index].Value
			fullIndex := lo.IndexOf(self.c.GetAppState().CustomCommandsHistory, item)
			if fullIndex == -1 {
				// Should never happen, but better be safe
				return nil
			}

			self.c.GetAppState().CustomCommandsHistory = slices.Delete(
				self.c.GetAppState().CustomCommandsHistory, fullIndex, fullIndex+1)
			self.c.SaveAppStateAndLogError()
			self.c.Contexts().Suggestions.RefreshSuggestions()
			return nil
		},
	})
}

func (self *CustomCommandAction) GetCustomCommandsHistorySuggestionsFunc() func(string) []*types.Suggestion {
	return func(input string) []*types.Suggestion {
		history := self.c.GetAppState().CustomCommandsHistory

		return helpers.FilterFunc(history, self.c.UserConfig.Gui.UseFuzzySearch())(input)
	}
}

// this mimics the shell functionality `ignorespace`
// which doesn't save a command to history if it starts with a space
func (self *CustomCommandAction) shouldSaveCommand(command string) bool {
	return !strings.HasPrefix(command, " ")
}
