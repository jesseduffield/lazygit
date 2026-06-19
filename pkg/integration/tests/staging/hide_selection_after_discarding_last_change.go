package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var HideSelectionAfterDiscardingLastChange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "After discarding the last change from the focused main view, the now-empty diff shows no selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = true
		config.GetUserConfig().Gui.SkipDiscardChangeWarning = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		// A single working-tree change, so discarding it empties the diff entirely.
		shell.UpdateFile("file1", "one\nTWO\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsShown().
			SelectedLines(
				Contains("-two"),
				Contains("+TWO"),
			).
			Press(keys.Universal.Remove)

		// The only change is gone, so the main view shows the placeholder with no
		// lingering selection.
		t.Views().Main().
			Content(Contains("No changed files")).
			SelectionIsHidden()
	},
})
