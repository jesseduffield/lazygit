package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageRangeFromMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select a range of diff lines in the focused main view and stage it, without diving into the staging view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "ONE\ntwo\nthree\nfour\nfive\nsix\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		// Focusing the main view selects the first change line (a single line).
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-one"),
			).
			// Select a range with `v` and extend it down past the modified line and
			// some context, then stage it.
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+four")).
			SelectedLines(
				Contains("-one"),
				Contains("+ONE"),
				Contains(" two"),
				Contains(" three"),
				Contains("+four"),
			).
			PressPrimaryAction().
			Tap(func() {
				// The selected change lines — the deletion, the addition replacing it,
				// and the added 'four' — are now staged; the context lines came along
				// but aren't themselves staged.
				t.Views().Secondary().
					ContainsLines(
						Contains("-one"),
						Contains("+ONE"),
						Contains(" two"),
						Contains(" three"),
						Contains("+four"),
					)
			}).
			// The unstaged half holds only the additions we didn't select.
			ContainsLines(
				Contains("+five"),
				Contains("+six"),
			)
	},
})
