package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ConfirmOnQuit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Quitting with a confirm prompt",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.ConfirmOnQuit = true
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Quit)

		t.ExpectPopup().Confirmation().
			Title(Equals("")).
			Content(Contains("Are you sure you want to quit?")).
			Confirm()
	},
})
