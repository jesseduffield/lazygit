package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FindBaseCommitForFixupScrollsIntoView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Finding the base commit for a fixup scrolls it into view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch").
			EmptyCommit("1st commit").
			CreateFileAndAdd("file1", "line 1\nline 2\nline 3\n").
			Commit("base commit").
			CreateNCommits(40).
			UpdateFile("file1", "line 1\nline 2 changed\nline 3\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Press(keys.Files.FindBaseCommitForFixup)

		// The base commit is at the very bottom of the list, far below the
		// visible area
		t.Views().Commits().
			IsFocused().
			SelectedLine(Contains("base commit")).
			SelectedLineIsVisible()
	},
})
