package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Both selection changes are handled before the next layout. Selecting file_a asks
// for its staged changes in the lower pane, whose task is only created after the
// layout. Selecting file_b again leaves that pane with nothing to show, and empties it
// right away. The pane is hidden then, but when it is shown again, it shows what it
// holds until its next render replaces it.
var KeepAnEmptiedPaneEmpty = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Selecting a file and moving back off it in rapid succession leaves the pane that the file filled empty",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file_a", "one\n")
		shell.CreateFileAndAdd("file_b", "one\n")
		shell.Commit("one")

		shell.UpdateFileAndAdd("file_a", "STAGED\n")
		shell.UpdateFile("file_a", "UNSTAGED\n")
		shell.UpdateFile("file_b", "two\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			NavigateToLine(Contains("file_b"))

		t.Views().Secondary().
			IsInvisible()

		t.Views().Files().
			PressRapidly(keys.Universal.PrevItem, keys.Universal.NextItem).
			SelectedLine(Contains("file_b"))

		t.Views().Secondary().
			IsInvisible().
			Content(Equals(""))
	},
})
