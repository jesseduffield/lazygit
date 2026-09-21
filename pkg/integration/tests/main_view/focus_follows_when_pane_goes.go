package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FocusFollowsWhenPaneGoes = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging the last unstaged change takes the upper pane away, so the focus follows the lines into the lower one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nADDED\ntwo\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Title(Equals("Unstaged changes")).
			SelectedLines(
				Contains("+ADDED"),
			).
			PressPrimaryAction()

		// Nothing is unstaged any more, so that pane is gone and the line is in the one
		// below, where the focus and the selection now are.
		t.Views().Files().Lines(
			Contains("M  file1"),
		)
		t.Views().Main().IsInvisible()
		t.Views().Secondary().
			IsFocused().
			Title(Equals("Staged changes")).
			SelectedLines(
				Contains("+ADDED"),
			)
	},
})
