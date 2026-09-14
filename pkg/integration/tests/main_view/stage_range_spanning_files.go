package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageRangeSpanningFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage a range reaching from one file's diff into another's, in a directory's focused main view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("fileA", "a\n")
		shell.CreateFileAndAdd("fileB", "b\n")
		shell.CreateFileAndAdd("fileC", "c\n")
		shell.Commit("one")

		shell.UpdateFile("fileA", "a\nfromA\n")
		shell.UpdateFile("fileB", "b\nfromB\n")
		shell.UpdateFile("fileC", "c\nfromC\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// With the root of the tree selected, the main view shows all three files' diffs.
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("▼ /").IsSelected(),
				Contains(" M fileA"),
				Contains(" M fileB"),
				Contains(" M fileC"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+fromA"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+fromB")).
			PressPrimaryAction()

		// Both files the range reached into are staged, each by its own patch, and the
		// file below it is untouched.
		t.Views().Files().Lines(
			Contains("▼ /"),
			Contains("M  fileA"),
			Contains("M  fileB"),
			Contains(" M fileC"),
		)
	},
})
