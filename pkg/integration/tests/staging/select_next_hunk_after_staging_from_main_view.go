package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectNextHunkAfterStagingFromMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "After staging a hunk from the focused main view, the selection advances to the next hunk rather than getting lost",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		// Two separate change blocks, far enough apart to stay distinct hunks.
		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			PressPrimaryAction().
			Tap(func() {
				t.Views().Secondary().
					ContainsLines(
						Contains("-three"),
						Contains("+THREE"),
					)
			}).
			// The selection didn't get lost: it advanced to the next (and now only
			// remaining) hunk, which is what the unstaged half shows.
			SelectedLines(
				Contains("-nine"),
				Contains("+NINE"),
			)
	},
})
