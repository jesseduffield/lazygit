package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PullStackedBranchesWithLocalCommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull a stack that was rebased on the remote; a branch below the current one that has a commit of its own is not offered for updating",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		createStackRewrittenOnTheRemote(shell)

		shell.Checkout("branch2")
		shell.EmptyCommit("mine")
		shell.Checkout("branch3")
		shell.RunCommand([]string{"git", "rebase", "--onto", "branch2", "branch2~1"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("branch3 ↓3↑4"),
				Contains("branch1 ↓1↑1"),
				Contains("branch2 ↓2↑3"),
				Contains("master"),
			)

		t.Views().Files().IsFocused().Press(keys.Universal.Pull)

		t.ExpectPopup().Menu().
			Title(Equals("Pull")).
			Lines(
				Contains("The following branches stacked below 'branch3' have also changed on the remote:"),
				Equals(""),
				Equals("  branch1 ↓1↑1"),
				Equals(""),
				Contains("Pull all these branches in addition to the current one"),
				Contains("Pull only 'branch3'"),
				Contains("Cancel"),
			).
			Cancel()
	},
})
