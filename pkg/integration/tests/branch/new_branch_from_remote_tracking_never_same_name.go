package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NewBranchFromRemoteTrackingNeverSameName = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Set tracking information when creating a new branch from a remote branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.Git.AutomaticTrackingWhenCreatingNewBranchFromRemoteBranch = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("main")
		shell.EmptyCommitWithDate("commit", "2023-04-07 10:00:00")
		shell.NewBranchFrom("other_branch", "main")
		shell.CloneIntoRemote("origin")
		shell.Checkout("main")
		shell.RunCommand([]string{"git", "branch", "-D", "other_branch"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Remotes().
			Focus().
			Lines(
				Contains("origin").IsSelected(),
			).
			PressEnter()

		t.Views().RemoteBranches().
			IsFocused().
			Lines(
				Contains("main").IsSelected(),
				Contains("other_branch"),
			).
			SelectNextItem().
			Press(keys.Universal.New)

		t.ExpectPopup().Prompt().
			Title(Equals("New branch name (branch is off of 'origin/other_branch')")).
			Confirm()

		t.Views().Branches().
			Focus().
			Lines(
				Contains("other_branch").IsSelected(),
				Contains("main"),
			).
			Press(keys.Branches.SetUpstream)

		t.ExpectPopup().Menu().Title(Contains("Upstream")).Select(Contains("View divergence from upstream")).Confirm()

		t.ExpectToast(Equals("Disabled: The selected branch has no upstream (or the upstream is not stored locally)"))
	},
})
