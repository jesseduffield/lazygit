package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddFromRemoteBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a new local tracking branch and worktree from a remote branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.NewBranch("feature")
		shell.CloneIntoRemote("origin")
		shell.Checkout("master")
		// drop the local branch so only the remote one remains
		shell.RunCommand([]string{"git", "branch", "-D", "feature"})
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
			NavigateToLine(Contains("feature")).
			Press(keys.Universal.NewWorktree).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("New worktree")).
					Select(Contains("New local branch and worktree from 'origin/feature'")).
					Confirm()

				// the new branch name defaults to the remote branch name with
				// the remote stripped off
				t.ExpectPopup().Prompt().
					Title(Equals("New branch and worktree name")).
					InitialText(Equals("feature")).
					Confirm()

				t.ExpectPopup().Menu().
					Title(Equals("Worktree location")).
					Confirm()
			})

		// we've switched into the new worktree, on a local branch that tracks
		// the remote one (the ✓ confirms tracking is set up)
		t.Views().Branches().
			IsFocused().
			Lines(
				Contains("feature").Contains("✓").IsSelected(),
				Contains("master (worktree repo)"),
			)
	},
})
