package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisabledKeybindings = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Confirms you can disable keybindings by setting them to <disabled>",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.Keybinding.Universal.PrevItem = "<disabled>"
		config.UserConfig.Keybinding.Universal.NextItem = "<disabled>"
		config.UserConfig.Keybinding.Universal.NextTab = "<up>"
		config.UserConfig.Keybinding.Universal.PrevTab = "<down>"
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press("<up>")

		t.Views().Worktrees().IsFocused()
	},
})
