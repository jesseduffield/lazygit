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
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Model().CommitCount(3)

		input.Views().Commits().
			Focus().
			Lines(
				Contains("commit 3"),
				Contains("commit 2"),
				Contains("commit 1"),
			).
			SelectNextItem().
			Press(keys.Universal.New)

		branchName := "my-branch-name"
		input.ExpectPrompt().Title(Contains("New Branch Name")).Type(branchName).Confirm()

		input.Model().CurrentBranchName(branchName)

		input.Views().Commits().Lines(
			Contains("commit 2"),
			Contains("commit 1"),
		)
	},
})
