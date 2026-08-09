package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectNextLineAfterStagingLineFromMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "In line mode, staging the first line of a block advances the selection to the next line of the same block, not to the next block",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "1\n2\n3\n4\n5\n6\n7\n8\n9\n")
		shell.Commit("one")

		// A three-line added block near the top, and a separate change block further
		// down (so "next block" is a distinct, wrong target for line mode).
		shell.UpdateFile("file1", "1\n2\nA\nB\nC\n3\n4\n5\n6\nX\n8\n9\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// Default (line) mode: the first added line is selected.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+A"),
			).
			// Stage just that line.
			PressPrimaryAction().
			// The selection advances to the next line of the same block, not to the
			// distant change block below.
			SelectedLines(
				Contains("+B"),
			)
	},
})
