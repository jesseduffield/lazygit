package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var InitialOpen = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Confirms a popup appears on first opening Lazygit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.DisableStartupPopups = false
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectPopup().Confirmation().
			Title(Equals("")).
			Content(Contains("Thanks for using lazygit!")).
			Confirm()

		t.Views().Files().IsFocused()
	},
})
