package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenLinkFailure = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When opening links via the OS fails, show a dialog instead.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.OpenLink = "exit 42"
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Information().Click(0, 0)

		t.ExpectPopup().Confirmation().
			Title(Equals("Error")).
			Content(Equals("Failed to open URL https://github.com/sponsors/jesseduffield\n\nError: exit status 42")).
			Confirm()
	},
})
