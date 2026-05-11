package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NewBranchFromRemoteTrackingDifferentName = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Set tracking information when creating a new branch from a remote branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("commit")
		shell.NewBranch("other_branch")
		shell.CloneIntoRemote("origin")
		shell.Checkout("master")
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
				Contains("master").IsSelected(),
				Contains("other_branch"),
			).
			SelectNextItem().
			Press(keys.Universal.New)

		t.ExpectPopup().Prompt().
			Title(Equals("New branch name (branch is off of 'origin/other_branch')")).
			Clear().
			Type("different_name").
			Confirm()

		t.Views().Branches().
			Focus().
			Lines(
				Contains("different_name").DoesNotContain("✓"),
				Contains("master"),
			)
	},
})
