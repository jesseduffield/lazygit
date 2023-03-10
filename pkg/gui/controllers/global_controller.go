package controllers

import (
	"strings"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type GlobalController struct {
	baseController
	*controllerCommon
}

func NewGlobalController(
	common *controllerCommon,
) *GlobalController {
	return &GlobalController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *GlobalController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.ExecuteCustomCommand),
			Handler:     self.customCommand,
			Description: self.c.Tr.LcExecuteCustomCommand,
		},
	}
}

func (self *GlobalController) customCommand() error {
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

			err := self.c.SaveAppState()
			if err != nil {
				self.c.Log.Error(err)
			}

			self.c.LogAction(self.c.Tr.Actions.CustomCommand)
			return self.c.RunSubprocessAndRefresh(
				self.os.Cmd.NewShell(command),
			)
		},
	})
}

// this mimics the shell functionality `ignorespace`
// which doesn't save a command to history if it starts with a space
func (self *GlobalController) shouldSaveCommand(command string) bool {
	return !strings.HasPrefix(command, " ")
}

func (self *GlobalController) GetCustomCommandsHistorySuggestionsFunc() func(string) []*types.Suggestion {
	// reversing so that we display the latest command first
	history := slices.Reverse(self.c.GetAppState().CustomCommandsHistory)

	return helpers.FuzzySearchFunc(history)
}

func (self *GlobalController) Context() types.Context {
	return nil
}
