package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterMenuCancelFilterWithEscape = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filtering the keybindings menu, then pressing esc to turn off the filter",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo:    func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Filter("Ignore").
			Lines(
				// menu has filtered down to the one item that matches the filter
				Contains(`--- Local ---`),
				Contains(`Ignore`).IsSelected(),
			)

		// Escape should cancel the filter, not close the menu
		t.GlobalPress(keys.Universal.Return)
		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			LineCount(GreaterThan(1))

		// Another escape closes the menu
		t.GlobalPress(keys.Universal.Return)
		t.Views().Files().IsFocused()
	},
})
