package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddFromBranchDetached = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add a detached worktree via the branches view",
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
					Select(Contains(`Create worktree from mybranch (detached)`)).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New worktree path")).
					Type("../linked-worktree").
					Confirm()
			}).
			// confirm we're still focused on the branches view
			IsFocused().
			Lines(
				Contains("(no branch)").IsSelected(),
				Contains("mybranch (worktree)"),
			)

		t.Views().Status().
			Content(Contains("repo(linked-worktree)"))
	},
})
