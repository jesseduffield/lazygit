package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FocusReturnsWhenSplitCollapses = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Unstaging the last staged change from the secondary pane brings the focus back to the main one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		// One change on each side, so the diff is split.
		shell.UpdateFileAndAdd("file1", "one\nSTAGED\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.UpdateFile("file1", "one\nSTAGED\ntwo\nthree\nfour\nfive\nsix\nseven\nUNSTAGED\neight\nnine\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("MM file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("+STAGED"),
			).
			PressPrimaryAction()

		// With nothing staged left the diff isn't split any more, so the pane that was
		// showing the staged side is gone — and the focus is back on the main one, where
		// the change just taken out of the index now is.
		t.Views().Files().Lines(
			Contains(" M file1"),
		)
		t.Views().Main().
			IsFocused().
			Content(Contains("+STAGED")).
			Content(Contains("+UNSTAGED")).
			SelectedLines(
				Contains("+STAGED"),
			)
	},
})
