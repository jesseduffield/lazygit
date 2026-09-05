package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SearchFollowsTheSelection = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stepping to the next match in the focused main view carries on from the selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "")
		shell.Commit("one")

		shell.UpdateFile("file1", "NEEDLE a\ntwo\nNEEDLE b\nfour\nNEEDLE c\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			FilterOrSearch("NEEDLE").
			SelectedLines(Contains("+NEEDLE a")).
			Tap(func() {
				t.Views().Search().Content(Contains("matches for 'NEEDLE' (1 of 3)"))
			}).
			// Move the selection past the second match by hand.
			SelectNextItem().
			SelectNextItem().
			SelectNextItem().
			SelectedLines(Contains("+four")).
			Tap(func() {
				t.Views().Search().Content(Contains("matches for 'NEEDLE' (2 of 3)"))
			}).
			// So the next match is the one after where the selection is, not the one
			// after the match it was last on.
			Press(keys.Universal.NextMatch).
			SelectedLines(Contains("+NEEDLE c"))
	},
})
