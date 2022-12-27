package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NewBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Creating a new branch from a commit",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("commit 1").
			EmptyCommit("commit 2").
			EmptyCommit("commit 3")
	},
	Run: func(shell *Shell, t *TestDriver, keys config.KeybindingConfig) {
		t.Model().CommitCount(3)

		t.Views().Commits().
			Focus().
			SelectNextItem().
			Lines(
				Contains("commit 3"),
				Contains("commit 2").IsSelected(),
				Contains("commit 1"),
			).
			Press(keys.Universal.New).
			Tap(func() {
				branchName := "my-branch-name"
				t.ExpectPrompt().Title(Contains("New Branch Name")).Type(branchName).Confirm()

				t.Model().CurrentBranchName(branchName)
			}).
			Lines(
				Contains("commit 2"),
				Contains("commit 1"),
			)
	},
})
