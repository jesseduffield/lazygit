package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FindBaseCommitForFixupInMainBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Finds the base commit to create a fixup for when the commit is already merged into master",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("1st commit").
			CreateFileAndAdd("file1", "file1 content\n").
			Commit("2nd commit").
			NewBranch("mybranch").
			CreateFileAndAdd("file2", "file2 content\n").
			Commit("3rd commit").
			EmptyCommit("4th commit").
			UpdateFile("file1", "file1 changed content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("4th commit").IsSelected(),
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
			Content(Contains("already on the main branch")).
			Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("4th commit"),
				Contains("3rd commit"),
				Contains("2nd commit").IsSelected(),
				Contains("1st commit"),
			)
	},
})
