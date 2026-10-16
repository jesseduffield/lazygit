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

		t.ExpectPopup().Menu().
			// a string that filters everything out
			Filter("ljasldkjaslkdjalskdjalsdjaslkd")

		t.Views().Menu().
			IsFocused().
			IsEmpty().
			// space is filter text in this menu, so we confirm with enter
			Press(keys.Universal.ConfirmMenu).
			Tap(func() {
				t.ExpectToast(Equals("Disabled: No item selected"))
			}).
			// escape the filter
			PressEscape().
			// escape the view
			PressEscape()

		// back in the files view, selecting the non-existing menu item was a no-op
		t.Views().Files().
			IsFocused()
	},
})
