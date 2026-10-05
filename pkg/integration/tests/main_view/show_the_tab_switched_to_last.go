package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Both tab switches are handled before the next layout. Switching to the Files tab
// asks for the file's diff, whose task is only created after the layout; switching
// back to the Worktrees tab asks for the worktree's details, whose task is created
// right away. The main view has to show what was asked for last.
var ShowTheTabSwitchedToLast = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Switching tabs twice in rapid succession leaves the main view showing the second tab's content",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\n")
		shell.Commit("one")
		shell.UpdateFile("file1", "ONE\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("(main worktree)").IsSelected(),
			)

		t.Views().Main().
			Content(Contains("Path:"))

		t.Views().Worktrees().
			PressRapidly(keys.Universal.PrevTab, keys.Universal.NextTab).
			IsFocused()

		t.Views().Main().
			Content(Contains("Path:"))
	},
})
