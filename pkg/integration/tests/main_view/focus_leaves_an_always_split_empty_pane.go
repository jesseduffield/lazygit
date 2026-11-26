package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FocusLeavesAnAlwaysSplitEmptyPane = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Configured to always split the diff, the emptied pane stays but the focus still leaves it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		cfg.GetUserConfig().Gui.SplitDiff = "always"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("one")

		shell.UpdateFileAndAdd("file1", "one\nSTAGED\ntwo\nthree\nfour\nfive\n")
		shell.UpdateFile("file1", "one\nSTAGED\ntwo\nthree\nUNSTAGED\nfour\nfive\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("+STAGED"),
			)

		t.Shell().Commit("two")
		t.GlobalPress(keys.Universal.Refresh)

		// The staged side is empty now, but its pane is still shown because the split
		// is configured as permanent. There is nothing left to act on in it, so the
		// focus goes where there is.
		t.Views().Secondary().
			IsVisible().
			Content(DoesNotContain("STAGED"))
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+UNSTAGED"),
			)
	},
})
