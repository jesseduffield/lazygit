package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ConfirmOnQuit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Quitting with a confirm prompt",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.ConfirmOnQuit = true
		config.UserConfig.Gui.Theme.ActiveBorderColor = []string{"red"}
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.CommitCount(0)

		input.PressKeys(keys.Universal.Quit)
		assert.InConfirm()
		assert.CurrentViewContent(Contains("Are you sure you want to quit?"))
		input.Confirm()
	},
})
