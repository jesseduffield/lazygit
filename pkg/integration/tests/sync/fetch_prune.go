package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FetchPrune = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fetch from the remote with the 'prune' option set in the git config",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// This option makes it so that git checks for deleted branches in the remote
		// upon fetching.
		shell.SetConfig("fetch.prune", "true")

		shell.EmptyCommit("my commit message")

		shell.NewBranch("branch_to_remove")
		shell.Checkout("master")
		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")
		shell.SetBranchUpstream("branch_to_remove", "origin/branch_to_remove")

		// # unbenownst to our test repo we're removing the branch on the remote, so upon
		// # fetching with prune: true we expect git to realise the remote branch is gone
		shell.RemoveRemoteBranch("origin", "branch_to_remove")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("master"),
				Contains("branch_to_remove").DoesNotContain("upstream gone"),
			)

		t.Views().Files().
			IsFocused().
			Press(keys.Files.Fetch)

		t.Views().Branches().
			Lines(
				Contains("master"),
				Contains("branch_to_remove").Contains("upstream gone"),
			)
	},
})
