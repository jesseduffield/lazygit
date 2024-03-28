package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterMenuWithNoKeybindings = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filtering the keybindings menu so that only entries without keybinding are left",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.Keybinding.Universal.ToggleWhitespaceInDiffView = "<disabled>"
	},
	SetupRepo: func(shell *Shell) {
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Filter("whitespace").
			Lines(
				// menu has filtered down to the one item that matches the
				// filter, and it doesn't have a keybinding
				Equals("--- Global ---"),
				Equals("Toggle whitespace").IsSelected(),
			)
	},
})
