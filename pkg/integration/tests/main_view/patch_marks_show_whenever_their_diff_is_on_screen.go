package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PatchMarksShowWheneverTheirDiffIsOnScreen = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The marks over the lines in the custom patch are shown whenever the main view shows the diff the patch is built from, whether or not it has the focus, and not over the diff of another commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			PressPrimaryAction().
			MarkedLines(
				Contains("-two"),
			).
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().IsFocused()
		t.Views().Main().MarkedLines(
			Contains("-two"),
		)

		// Leaving the main view keeps them, as it keeps the patch previewed beside the
		// diff.
		t.Views().Secondary().Press(keys.Universal.Return)

		t.Views().Commits().IsFocused()
		t.Views().Main().MarkedLines(
			Contains("-two"),
		)

		// The diff of another commit has none of the patch's lines.
		t.Views().Commits().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			SelectNextItem()

		t.Views().Main().NoMarkedLines()

		// And they are back with the diff they belong to.
		t.Views().Commits().SelectPreviousItem()

		t.Views().Main().MarkedLines(
			Contains("-two"),
		)
	},
})
