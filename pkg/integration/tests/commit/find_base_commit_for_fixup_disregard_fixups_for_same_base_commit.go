package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FindBaseCommitForFixupDisregardFixupsForSameBaseCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Finds the base commit to create a fixup for, disregarding fixup commits for the same base commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("1st commit").
			NewBranch("mybranch").
			CreateFileAndAdd("file1", "1\n2\n3\n").
			Commit("2nd commit").
			UpdateFileAndAdd("file1", "1\n2\n3a\n").
			Commit("fixup! 2nd commit").
			EmptyCommit("3rd commit").
			UpdateFileAndAdd("file1", "1a\n2\n3b\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("3rd commit").IsSelected(),
				Contains("fixup! 2nd commit"),
				Contains("2nd commit"),
				Contains("1st commit"),
			)

		t.Views().Files().
			Focus().
			Press(keys.Files.FindBaseCommitForFixup)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(
				MatchesRegexp("Multiple base commits found.*\n\n" +
					".*fixup! 2nd commit\n" +
					".*2nd commit"),
			).
			Confirm()
	},
})
