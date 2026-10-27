package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NoSelectionOverACommitLog = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focusing the main view over a branch's commit log shows no selection, even for a log too long to be read in one go",
	ExtraCmdArgs: []string{},
	Skip:         false,
	// A short terminal, so that the log below is longer than the initial read of it.
	Width:       100,
	Height:      20,
	SetupConfig: func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(60)
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("add file1")
		shell.UpdateFile("file1", "one\ntwo modified\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Leave a selection behind in the main view, on the file's diff.
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsActive().
			PressEscape()

		t.Views().Branches().
			Focus()

		t.Views().Main().
			Content(Contains("commit-60"))

		// A commit log holds nothing to point at, so focusing it shows no selection —
		// not even the one the pane was left with under the files panel.
		t.Views().Branches().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsHidden()
	},
})
