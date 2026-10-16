package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectHunkBelowLastChange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggling hunk selection while below the last change selects the last change block",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nTWO\nthree\nfour\nfive\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// Below the last change there is no block ahead to select, so hunk mode takes
		// the one behind rather than doing nothing.
		t.Views().Main().
			IsFocused().
			NavigateToLine(Contains(" five")).
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-two"),
				Contains("+TWO"),
			)
	},
})
