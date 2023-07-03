package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filtering the keybindings menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains(`??`).Contains(`myfile`).IsSelected(),
			).
			Press(keys.Universal.OptionMenu).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Keybindings")).
					Filter("Toggle staged").
					Lines(
						// menu has filtered down to the one item that matches the filter
						Contains(`Toggle staged`).IsSelected(),
					).
					Confirm()
			})

		t.Views().Files().
			IsFocused().
			Lines(
				// file has been staged
				Contains(`A `).Contains(`myfile`).IsSelected(),
			).
			// Upon opening the menu again, the filter should have been reset
			Press(keys.Universal.OptionMenu).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Keybindings")).
					LineCount(GreaterThan(1))
			})
	},
})
