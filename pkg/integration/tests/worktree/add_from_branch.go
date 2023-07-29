package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddFromBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add a worktree via the branches view, then switch back to the main worktree via the branches view",
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
					Title(Equals("Worktree")).
					Select(Contains(`Create worktree from mybranch`).DoesNotContain("detached")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New worktree path")).
					Type("../linked-worktree").
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New branch name")).
					Type("newbranch").
					Confirm()
			}).
			// confirm we're still focused on the branches view
			IsFocused().
			Lines(
				Contains("newbranch").IsSelected(),
				Contains("mybranch (worktree)"),
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
				Contains("newbranch (worktree)"),
			).
			// Confirm the files view is still showing in the files window
			Press(keys.Universal.PrevBlock)

		t.Views().Files().
			IsFocused()
	},
})
