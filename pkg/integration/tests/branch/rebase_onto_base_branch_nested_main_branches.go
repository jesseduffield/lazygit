package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebaseOntoBaseBranchNestedMainBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When several nested main branches contain the fork point, the first configured one is the base branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.MainBranches = []string{"master", "rc", "dev"}
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("master 1").
			NewBranch("rc").
			EmptyCommit("rc 1").
			NewBranch("dev").
			EmptyCommit("dev 1").
			NewBranchFrom("feature", "master").
			EmptyCommit("feature 1").
			CloneIntoRemote("origin").
			SetBranchUpstream("master", "origin/master").
			SetBranchUpstream("rc", "origin/rc").
			SetBranchUpstream("dev", "origin/dev")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("feature").IsSelected(),
				Contains("dev"),
				Contains("rc"),
				Contains("master"),
			).
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Rebase 'feature'")).
			ContainsLines(Contains("Rebase onto base branch (origin/master)"))
	},
})
