package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PullStackedBranchesCheckedOutInOtherWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull a stack that was rebased on the remote; a branch below the current one that is checked out in another worktree is not offered for updating",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		createStackRewrittenOnTheRemote(shell)
		shell.Checkout("branch3")
		shell.SetConfig("pull.rebase", "true")

		shell.AddWorktreeCheckout("branch1", "../linked-worktree")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().Press(keys.Universal.Pull)

		t.ExpectPopup().Menu().
			Title(Equals("Pull")).
			Lines(
				Contains("The following branches stacked below 'branch3' have also changed on the remote:"),
				Equals(""),
				Equals("  branch2 ↓2↑2"),
				Equals(""),
				Contains("Pull all these branches in addition to the current one"),
				Contains("Pull only 'branch3'"),
				Contains("Cancel"),
			).
			Select(Contains("Pull all these branches in addition to the current one")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("branch3 ✓"),
				Contains("branch1 (worktree linked-worktree) ↓1↑1"),
				Contains("branch2 ✓"),
				Contains("master"),
			)
	},
})
