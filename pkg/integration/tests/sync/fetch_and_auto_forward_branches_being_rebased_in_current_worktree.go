package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FetchAndAutoForwardBranchesBeingRebasedInCurrentWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Auto-forward skips a main branch that is being rebased in the current worktree",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.AutoForwardBranches = "onlyMainBranches"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(3)
		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")
		shell.HardReset("HEAD^")

		// the failing exec stops the rebase after picking commit-02, with HEAD
		// detached from master
		shell.RunCommandExpectError([]string{"git", "rebase", "--exec", "false", "HEAD^"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("(no branch, rebasing master)"),
				Contains("master ↓1"),
			)

		t.Views().Files().
			IsFocused().
			Press(keys.Files.Fetch)

		t.Views().Branches().
			Lines(
				Contains("(no branch, rebasing master)"),
				Contains("master ↓1"),
			)
	},
})
