package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FindBaseCommitForFixupIgnoreFixupCommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Finds the base commit to create a fixup for, disregarding changes to a commit that is already on master",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("1st commit").
			CreateFileAndAdd("file", "file1 content line 1\nfile1 content line 2\n").
			Commit("2nd commit").
			NewBranch("mybranch").
			UpdateFileAndAdd("file", "file1 changed content line 1\nfile1 changed content line 2\n").
			Commit("3rd commit").
			EmptyCommit("4th commit").
			UpdateFileAndAdd("file", "file1 1st fixup content line 1\nfile1 changed content line 2\n").
			Commit("fixup! 3rd commit").
			UpdateFile("file", "file1 2nd fixup content line 1\nfile1 2nd fixup content line 2\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("fixup! 3rd commit").IsSelected(),
				Contains("4th commit"),
				Contains("3rd commit"),
				Contains("2nd commit"),
				Contains("1st commit"),
			)

		t.Views().Files().
			Focus().
			Press(keys.Files.FindBaseCommitForFixup)
		t.ExpectPopup().
			Confirmation().
			Title(Equals("Find base commit for fixup")).
			Content(
				Contains("all but one of them were fixup commits").
					Contains("3rd commit").
					Contains("fixup! 3rd commit"),
			).
			Confirm()
		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("fixup! 3rd commit"),
				Contains("4th commit"),
				Contains("3rd commit").IsSelected(),
				Contains("2nd commit"),
				Contains("1st commit"),
			)
	},
})
