package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PatchMarksShowWhileTheDiffIsFocused = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The marks over the lines in the custom patch are shown while either pane of the focused main view holds the focus, and not once it leaves",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
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
			// Moving to the pane previewing the patch is still working on the same patch,
			// so the marks stay.
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().IsFocused()
		t.Views().Main().MarkedLines(
			Contains("-two"),
		)

		// Leaving the diff behind takes them away: they say what pressing space here
		// would act on.
		t.Views().Secondary().Press(keys.Universal.Return)

		t.Views().Commits().IsFocused()
		t.Views().Main().NoMarkedLines()

		// And they are back with the focus.
		t.Views().Commits().Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			MarkedLines(
				Contains("-two"),
			)
	},
})
