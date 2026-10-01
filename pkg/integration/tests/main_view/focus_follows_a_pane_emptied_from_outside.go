package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FocusFollowsAPaneEmptiedFromOutside = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Committing the staged changes outside lazygit takes the staged pane away, so the focus follows into the one that is left",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("one")

		// One change on each side, so the diff is split.
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
			Title(Equals("Staged changes")).
			SelectedLines(
				Contains("+STAGED"),
			)

		// Nothing lazygit did empties the staged side here; the refresh simply finds
		// it empty, and the pane the focus was in is gone by the time it lands.
		t.Shell().Commit("two")
		t.GlobalPress(keys.Universal.Refresh)

		t.Views().Secondary().IsInvisible()
		t.Views().Main().
			IsFocused().
			Title(Equals("Unstaged changes")).
			SelectedLines(
				Contains("+UNSTAGED"),
			)
	},
})
