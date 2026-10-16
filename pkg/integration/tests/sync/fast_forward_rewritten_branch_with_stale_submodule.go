package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardRewrittenBranchWithStaleSubmodule = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fast-forward a branch with a rewritten upstream branch while its worktree has a submodule checked out at a different commit than the one recorded",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial")
		shell.CloneIntoSubmodule("sub", "sub")
		createBranchRewrittenOnTheRemote(shell)

		shell.Checkout("feature")
		shell.RunCommand([]string{"git", "-C", "sub", "commit", "--allow-empty", "-m", "newer"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Lines(
				Contains("sub"),
			)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("feature ↓2↑2").IsSelected(),
				Contains("master"),
			).
			Press(keys.Branches.FastForward).
			Lines(
				Contains("feature ✓").IsSelected(),
				Contains("master"),
			)

		// The submodule was left alone
		t.Views().Files().
			Lines(
				Contains("sub"),
			)
	},
})
