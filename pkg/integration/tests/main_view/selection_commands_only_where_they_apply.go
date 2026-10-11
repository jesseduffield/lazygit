package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectionCommandsOnlyWhereTheyApply = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The keybindings menu offers the commands that act on a diff selection only in a main view that has one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// A branch's commit log has nothing to select, so the commands that act on a
		// selection have no business being listed there.
		t.Views().Branches().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Tap(func() {
				t.Views().Menu().
					Content(DoesNotContain("Select hunks")).
					Content(DoesNotContain("Toggle range select")).
					Content(DoesNotContain("Go to next hunk")).
					Content(DoesNotContain("Go to next file"))
			}).
			Cancel()

		// A diff view does list them, and with nothing changed to select they're
		// offered but disabled.
		t.Views().Files().
			Focus().
			IsEmpty().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsHidden().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Select(Contains("Select hunks")).
			Confirm()

		t.ExpectToast(Contains("There is nothing to select here"))
	},
})
