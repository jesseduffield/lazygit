package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveCommitsToNewBranchNoUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Moving commits to a new branch is not allowed when the current branch has no upstream branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.EmptyCommit("new commit 1")
		shell.EmptyCommit("new commit 2")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
			).
			Press(keys.Branches.MoveCommitsToNewBranch)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Cannot move commits from a branch that has no upstream branch")).
			Confirm()
	},
})
