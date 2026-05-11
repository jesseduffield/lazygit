package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EmptyMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that we don't crash on an empty menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		t.Views().Menu().
			IsFocused().
			// a string that filters everything out
			FilterOrSearch("ljasldkjaslkdjalskdjalsdjaslkd").
			IsEmpty().
			Press(keys.Universal.Select).
			Tap(func() {
				t.ExpectToast(Equals("Disabled: No item selected"))
			}).
			// escape the search
			PressEscape().
			// escape the view
			PressEscape()

		// back in the files view, selecting the non-existing menu item was a no-op
		t.Views().Files().
			IsFocused()
	},
})
