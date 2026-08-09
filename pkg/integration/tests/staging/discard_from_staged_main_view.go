package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardFromStagedMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding a hunk on the staged half of the focused main view just unstages it, advancing the selection to the next staged hunk like unstaging does",
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
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			// Discarding on the staged side is just unstaging — no discard confirm — so the
			// selection stays in the staged half and advances to the next staged hunk
			// rather than getting lost or jumping to the unstaged half.
			Press(keys.Universal.Remove)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("-nine"),
				Contains("+NINE"),
			)

		// The unstaged change is untouched.
		t.Views().Main().
			Content(Contains("-six")).
			Content(Contains("+SIX"))
	},
})
