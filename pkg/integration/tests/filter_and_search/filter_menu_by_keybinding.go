package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterMenuByKeybinding = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filtering the keybindings menu by keybinding",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Press(keys.Universal.OptionMenu).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Keybindings")).
					Filter("@+").
					Lines(
						// menu has filtered down to the one item that matches the filter
						Contains("--- Global ---"),
						Contains("+ Next screen mode").IsSelected(),
					).
					Confirm()
			}).

			// Upon opening the menu again, the filter should have been reset
			Press(keys.Universal.OptionMenu).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Keybindings")).
					LineCount(GreaterThan(1))
			})
	},
})
