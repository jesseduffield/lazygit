package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectNextChangeAfterUnstaging = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Taking one line of a staged modification back out of the index moves on to the line that replaced it, in the pane the staged side ended up in",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "c0\na1\na2\na3\nc1\nc2\noldB\nc3\n")
		shell.Commit("one")

		// Staged: a block of deletions, and a modification below it.
		shell.UpdateFileAndAdd("file1", "c0\nc1\nc2\nnewB\nc3\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// The file has nothing but staged changes, so the pane holding them is the only
		// one shown.
		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("-a1"),
			).
			NavigateToLine(Contains("-oldB")).
			PressPrimaryAction()

		// The unstaged pane has appeared above, and the work carries on where it was:
		// on the line that takes the place of the one taken out.
		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("+newB"),
			)
	},
})
