package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FocusFollowsStagedSide = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Unstaging from a fully staged file leaves the focus on the staged side, which keeps its pane",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		// Two staged additions and nothing unstaged, so only the pane the staged side
		// lives in is shown.
		shell.UpdateFileAndAdd("file1", "one\nSTAGED1\ntwo\nthree\nfour\nfive\nsix\nseven\nSTAGED2\neight\nnine\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("M  file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().IsInvisible()
		t.Views().Secondary().
			IsFocused().
			Title(Equals("Staged changes")).
			SelectedLines(
				Contains("+STAGED1"),
			).
			PressPrimaryAction()

		// The line taken out of the index turns up in the pane that has just appeared
		// above, and the work carries on where it was, on the next staged change.
		t.Views().Files().Lines(
			Contains("MM file1"),
		)
		t.Views().Main().
			IsVisible().
			Content(Contains("+STAGED1"))
		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("+STAGED2"),
			)
	},
})
