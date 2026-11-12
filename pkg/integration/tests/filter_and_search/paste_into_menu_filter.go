package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PasteIntoMenuFilter = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Paste text into the filter of a menu that filters as you type",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Paste("\nIgnore\n").
			Lines(
				Contains(`─── Local`),
				Contains(`Ignore`).IsSelected(),
			)

		t.Views().MenuFilter().
			Content(Equals("Ignore"))
	},
})
