package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AdvanceAfterStagingShiftsLineNumbers = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a hunk that adds a line still leaves the selection on the next hunk, whose line numbers it moved",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "1\n2\n3\n4\n5\n6\n7\n8\n")
		shell.Commit("one")

		// Three change blocks: a modification, an added line, and another modification
		// below it. Staging the middle one changes how many lines the file has, and so
		// where the last one sits.
		shell.UpdateFile("file1", "1\nX\n3\n4\nNEW\n5\n6\nY\n8\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-2"),
				Contains("+X"),
			).
			Press(keys.Main.NextHunk).
			SelectedLines(
				Contains("+NEW"),
			).
			PressPrimaryAction().
			SelectedLines(
				Contains("-7"),
				Contains("+Y"),
			)
	},
})
