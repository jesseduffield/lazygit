package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterMenuKeyHandling = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Which keys drive a menu that filters as you type, and which ones are filter text",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// so that quitting is observable instead of ending the test
		cfg.GetUserConfig().ConfirmOnQuit = true
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// presses a key that is expected to move the selection away from the first
		// item, and one that is expected to bring it back
		navigates := func(forward string, back string) {
			t.GlobalPress(config.Keybinding{forward})
			t.Views().Menu().SelectedLineIdxAtLeast(2)
			t.GlobalPress(config.Keybinding{back})
			t.Views().Menu().SelectedLineIdx(1)
		}

		t.Views().Files().IsFocused().Press(keys.Universal.OptionMenu)
		t.Views().Menu().IsFocused().SelectedLineIdx(1)

		// Until there is a filter, the configured navigation keys drive the menu,
		// printable or not
		navigates("<down>", "<up>")
		navigates("j", "k")
		navigates(".", ",")
		navigates(">", "<")
		t.Views().MenuFilter().IsInvisible()

		// A menu item's own key is filter text; it doesn't execute the item. 'c'
		// commits when the files view has the focus.
		t.ExpectPopup().Menu().Filter("c")
		t.Views().Menu().IsFocused()
		t.Views().MenuFilter().Content(Equals("c"))

		// So is the key that filters other lists
		t.GlobalPress(keys.Universal.StartSearch)
		t.Views().MenuFilter().Content(Equals("c/"))
		t.Views().Search().IsInvisible()

		// And so are the printable navigation keys, now that there is somewhere for
		// them to go
		t.GlobalPress(config.Keybinding{"<c-u>"})
		t.GlobalPress(config.Keybinding{"j"})
		t.GlobalPress(config.Keybinding{"."})
		t.GlobalPress(config.Keybinding{">"})
		t.Views().MenuFilter().Content(Equals("j.>"))

		// The keys that can't be typed keep driving the menu
		t.GlobalPress(config.Keybinding{"<c-u>"})
		navigates("<down>", "<up>")
		navigates("<pgdown>", "<pgup>")
		navigates("<end>", "<home>")

		// Keys that the filter doesn't take and the menu doesn't handle reach the
		// global keybindings
		t.GlobalPress(config.Keybinding{"<c-c>"})
		t.ExpectPopup().Confirmation().
			Title(Equals("")).
			Content(Contains("Are you sure you want to quit?")).
			Confirm()
	},
})
