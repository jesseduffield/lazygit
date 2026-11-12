package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PasteIntoKeybindingsMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pasting text that ends with a newline into the keybindings menu filters the menu without running the selected item",
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
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Paste("stage all\n")

		t.Views().Files().
			Lines(
				Contains(`??`).Contains(`myfile`),
			)
	},
})
