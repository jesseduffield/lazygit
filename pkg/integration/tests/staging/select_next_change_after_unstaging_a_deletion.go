package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectNextChangeAfterUnstagingADeletion = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Unstaging the deletion of a modification advances to its replacement line, not back to an earlier block (the staged diff's new side is the index, which unstaging shifts)",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "c0\na1\na2\na3\nc1\nc2\noldB\nc3\n")
		shell.Commit("one")
		// Staged: a block of deletions (a1-a3) and, below it, a modification
		// (oldB -> newB).
		shell.UpdateFileAndAdd("file1", "c0\nc1\nc2\nnewB\nc3\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().Press(keys.Universal.FocusMainView)

		// Only-staged file: the main view shows the staged diff; navigate down to the
		// deletion of the modification.
		t.Views().Main().IsFocused().
			SelectedLines(Contains("-a1")).
			Press(keys.Universal.NextItem).
			Press(keys.Universal.NextItem).
			Press(keys.Universal.NextItem).
			Press(keys.Universal.NextItem).
			Press(keys.Universal.NextItem).
			SelectedLines(Contains("-oldB")).
			PressPrimaryAction()

		// Unstaging splits the file; focus follows the staged remainder into the
		// secondary half and lands on the modification's replacement line.
		t.Views().Secondary().IsFocused().
			SelectedLines(Contains("+newB"))
	},
})
