package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectNextHunkAfterUnstagingFromMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "After unstaging a hunk from the staged half of the focused main view, the selection advances to the next staged hunk",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\neleven\n")
		shell.Commit("one")

		// Two staged hunks...
		shell.UpdateFileAndAdd("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\neleven\n")
		// ...plus an unstaged change, so the main view splits into staged/unstaged.
		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\nSIX\nseven\neight\nNINE\nten\neleven\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// The unstaged half is focused first; switch to the staged half.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-six"),
				Contains("+SIX"),
			).
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			// Unstage the first staged hunk.
			PressPrimaryAction().
			// The selection advances to the next staged hunk rather than getting lost.
			SelectedLines(
				Contains("-nine"),
				Contains("+NINE"),
			)
	},
})
