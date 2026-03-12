package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoveWorktreeFromBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a worktree and delete its branch from the branches view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.EmptyCommit("commit 2")
		shell.EmptyCommit("commit 3")
		shell.AddWorktree("mybranch", "../linked-worktree", "newbranch")
		shell.AddFileInWorktreeOrSubmodule("../linked-worktree", "file", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("mybranch").IsSelected(),
				Contains("newbranch (worktree linked-worktree)"),
			).
			NavigateToLine(Contains("newbranch")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Delete branch 'newbranch'?")).
					Select(Contains("Remove worktree")).
					Confirm()

				t.ExpectPopup().Confirmation().
					Title(Equals("Remove worktree")).
					Content(Equals("Are you sure you want to remove worktree 'linked-worktree'?")).
					Confirm()

				t.ExpectPopup().Confirmation().
					Title(Equals("Remove worktree")).
					Content(Equals("'linked-worktree' contains modified or untracked files, or submodules (or all of these). Are you sure you want to remove it?")).
					Confirm()
			}).
			Lines(
				Contains("mybranch").IsSelected(),
			)

		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("(main worktree)").IsSelected(),
			)
	},
})
