package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FindBaseCommitForFixupWarningForAddedLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Finds the base commit to create a fixup for, and warns that there are hunks with only added lines",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch").
			EmptyCommit("1st commit").
			CreateFileAndAdd("file1", "file1 content\n").
			Commit("2nd commit").
			CreateFileAndAdd("file2", "file2 content\n").
			Commit("3rd commit").
			UpdateFile("file1", "file1 changed content").
			UpdateFile("file2", "file2 content\nadded content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("3rd commit").IsSelected(),
				Contains("2nd commit"),
				Contains("1st commit"),
			)

		t.Views().Files().
			Focus().
			Press(keys.Files.FindBaseCommitForFixup)

		t.ExpectPopup().Confirmation().
			Title(Equals("Find base commit for fixup")).
			Content(Contains("There are ranges of only added lines in the diff; be careful to check that these belong in the found base commit.")).
			Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("3rd commit"),
				Contains("2nd commit").IsSelected(),
				Contains("1st commit"),
			)
	},
})
