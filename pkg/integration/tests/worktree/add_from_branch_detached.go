package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddFromBranchDetached = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add a detached worktree at a branch via the branches view, choosing a custom location",
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
			Press(keys.Universal.NewWorktree).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("New worktree")).
					Select(Contains("New detached worktree at 'mybranch'")).
					Confirm()

				// the location menu defaults the directory name to the branch
				// name; pick "Other…" to type a different path instead
				t.ExpectPopup().Menu().
					Title(Equals("Worktree location")).
					Select(Contains("Other…")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New worktree path")).
					InitialText(Contains("mybranch")).
					Clear().
					Type("../linked-worktree").
					Confirm()
			}).
			// confirm we're still focused on the branches view
			IsFocused().
			Lines(
				Contains("(no branch)").IsSelected(),
				Contains("mybranch (worktree repo)"),
			)

		t.Views().Status().
			Content(Contains("repo(linked-worktree)"))
	},
})
