package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddFromBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a new branch and worktree from a branch, then switch back to the main worktree via the branches view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("mybranch"),
			).
			Press(keys.Worktrees.ViewWorktreeOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("New worktree")).
					Select(Contains("New branch and worktree from 'mybranch'")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New branch and worktree name")).
					Type("newbranch").
					Confirm()

				// no existing worktrees and no configured default path, so the
				// only candidate location is the repo's parent directory; accept it
				t.ExpectPopup().Menu().
					Title(Equals("Worktree location")).
					Confirm()
			}).
			// confirm we're still focused on the branches view
			IsFocused().
			Lines(
				Contains("newbranch").IsSelected(),
				Contains("mybranch (worktree repo)"),
			).
			NavigateToLine(Contains("mybranch")).
			Press(keys.Universal.Select).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Switch to worktree")).
					Content(Equals("This branch is checked out by worktree repo. Do you want to switch to that worktree?")).
					Confirm()
			}).
			Lines(
				Contains("mybranch").IsSelected(),
				// the worktree's directory name matches the branch name, so the
				// branches view shows the compact "(worktree)" with no name
				Contains("newbranch (worktree)"),
			).
			// Confirm the files view is still showing in the files window
			Press(keys.Universal.PrevBlock)

		t.Views().Files().
			IsFocused()
	},
})
