package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UnstageDiffLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Take a line back out of the index from the staged half of the focused main view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		// Two staged additions, far enough apart to be separate hunks, plus an unstaged
		// one, so that the diff is split into a staged and an unstaged half.
		shell.UpdateFileAndAdd("file1", "one\nSTAGED1\ntwo\nthree\nfour\nfive\nsix\nseven\nSTAGED2\neight\nnine\nten\n")
		shell.UpdateFile("file1", "one\nSTAGED1\ntwo\nthree\nUNSTAGED\nfour\nfive\nsix\nseven\nSTAGED2\neight\nnine\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("MM file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		// The main half holds the unstaged changes; the staged ones are next door.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+UNSTAGED"),
			).
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("+STAGED1"),
			).
			PressPrimaryAction()

		// The line acted on is out of the index, and the one below it stays in.
		t.Views().Secondary().
			Content(DoesNotContain("+STAGED1")).
			ContainsLines(
				Contains("+STAGED2"),
			)
		t.Views().Main().ContainsLines(
			Contains("+STAGED1"),
		)
	},
})
