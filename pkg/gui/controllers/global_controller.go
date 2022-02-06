package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
			self.c.GetAppState().CustomCommandsHistory = utils.Limit(
				utils.Uniq(
					append(self.c.GetAppState().CustomCommandsHistory, command),
				),
				1000,
			)

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

func (self *GlobalController) GetCustomCommandsHistorySuggestionsFunc() func(string) []*types.Suggestion {
	// reversing so that we display the latest command first
	history := utils.Reverse(self.c.GetAppState().CustomCommandsHistory)

	return helpers.FuzzySearchFunc(history)
}

func (self *GlobalController) Context() types.Context {
	return nil
}
