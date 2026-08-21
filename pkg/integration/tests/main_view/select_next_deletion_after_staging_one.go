package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectNextDeletionAfterStagingOne = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging one deletion in the middle of a block of them moves on to the next, not back to the block's first line",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "keep1\nd1\nd2\nd3\nd4\nkeep2\n")
		shell.Commit("one")

		// Four deletions in a row. They all sit at the same place in the new file, which
		// is what makes them impossible to tell apart by that alone.
		shell.UpdateFile("file1", "keep1\nkeep2\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-d1"),
			).
			Press(keys.Universal.NextItem).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("-d3"),
			).
			PressPrimaryAction().
			SelectedLines(
				Contains("-d4"),
			)
	},
})
