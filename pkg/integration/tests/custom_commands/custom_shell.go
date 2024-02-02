package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CustomShell = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Confirms a popup appears on first opening Lazygit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.OS.Shell = "sh"
		config.UserConfig.OS.ShellArg = "-c"
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsEmpty().
			IsFocused().
			Press(keys.Universal.ExecuteCustomCommand)

		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			Type("echo hello world").
			Confirm()
	},
})
