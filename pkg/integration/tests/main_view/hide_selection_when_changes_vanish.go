package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var HideSelectionWhenChangesVanish = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The main view's selection disappears along with the changes it was on",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsActive().
			Press(keys.Universal.Return)

		// Discarding the change leaves the main view with a placeholder to show, so the
		// selection that was on the change goes with it rather than lingering over the
		// message.
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			IsEmpty()

		t.Views().Main().
			Content(Contains("No changed files")).
			SelectionIsHidden()
	},
})
