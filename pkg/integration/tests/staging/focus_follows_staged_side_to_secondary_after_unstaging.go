package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FocusFollowsStagedSideToSecondaryAfterUnstaging = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Unstaging the first hunk of an only-staged file from the focused main view splits the diff; focus follows the staged remainder into the secondary half",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\neleven\n")
		shell.Commit("one")

		// Two staged hunks and no unstaged changes, so the main view shows the staged
		// diff in the main half without a split.
		shell.UpdateFileAndAdd("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\neleven\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// The only-staged file shows the staged diff in the main half, first hunk selected.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			// Unstage the first staged hunk, which splits the file into staged + unstaged.
			PressPrimaryAction()

		// Focus follows the staged remainder into the secondary half, landing on the
		// next staged hunk rather than staying on the now-unstaged half.
		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("-nine"),
				Contains("+NINE"),
			)

		// The main half now shows the hunk we just unstaged.
		t.Views().Main().
			ContainsLines(
				Contains("-three"),
				Contains("+THREE"),
			)
	},
})
