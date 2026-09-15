package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeybindingSuggestionsDontCrashOnDisabledBindings = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter out keybinding suggestions whose bindings are disabled",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Keybinding.Files.StashAllChanges = []string{}
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus()
		t.Views().Options().Content(
			Equals("Commit: c | Reset: D | Keybindings: ?"))
	},
})
