package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var JumpToAFileOfAStashWithNoPrefixes = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Jump to a file in a directory named b of a stash entry's diff when diff.noprefix is set",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("b/file1", "one\n")
		shell.CreateFileAndAdd("file2", "one\n")
		shell.Commit("one")

		shell.UpdateFile("b/file1", "two\n")
		shell.UpdateFile("file2", "two\n")
		shell.Stash("stash one")

		// git then leaves out the a/ and b/ in front of the paths of a diff. The b/
		// that the paths of file1 start with is its directory.
		shell.SetConfig("diff.noprefix", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().
			Focus().
			Lines(
				Contains("stash one").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsActive().
			Press(keys.Universal.JumpToFile)

		t.ExpectPopup().Menu().
			Title(Equals("Jump to file")).
			Lines(
				/* EXPECTED:
				Equals("b/file1"),
				ACTUAL: */
				Equals("file1"),
				Equals("file2"),
				Equals("Cancel"),
			)
	},
})
