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
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.CommitCount(0)

		input.Press(keys.Universal.Quit)
		input.Confirmation().
			Title(Equals("")).
			Content(Contains("Are you sure you want to quit?")).
			Confirm()
	},
})
