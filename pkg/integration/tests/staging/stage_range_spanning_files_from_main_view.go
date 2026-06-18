package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageRangeSpanningFilesFromMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select a range spanning two files in a directory's focused main view and stage it in one go",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("a", "a\n")
		shell.CreateFileAndAdd("b", "b\n")
		shell.Commit("one")

		shell.UpdateFile("a", "a\nfromA\n")
		shell.UpdateFile("b", "b\nfromB\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// The root node is selected, so the focused main view shows both files' diffs.
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("▼ /").IsSelected(),
				Contains(" M a"),
				Contains(" M b"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+fromA"),
			).
			// Select a range reaching from the addition in the first file into the
			// addition in the second, then stage it.
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+fromB")).
			PressPrimaryAction().
			Tap(func() {
				// Both files' additions got staged in one go.
				t.Views().Files().Lines(
					Contains("▼ /"),
					Contains("M  a"),
					Contains("M  b"),
				)
			})
	},
})
