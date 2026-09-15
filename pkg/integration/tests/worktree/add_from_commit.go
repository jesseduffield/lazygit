package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddFromCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a new branch and worktree from a commit via the commits view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.EmptyCommit("commit two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit two").IsSelected(),
				Contains("initial commit"),
			).
			NavigateToLine(Contains("initial commit")).
			Press(keys.Universal.NewWorktree).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("New worktree")).
					Select(Contains("New branch and worktree from")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New branch and worktree name")).
					Type("newbranch").
					Confirm()

				t.ExpectPopup().Menu().
					Title(Equals("Worktree location")).
					Confirm()
			}).
			Lines(
				Contains("initial commit"),
			)

		// Confirm we're now in the branches view
		t.Views().Branches().
			IsFocused().
			Lines(
				Contains("newbranch").IsSelected(),
				Contains("mybranch (worktree repo)"),
			)
	},
})
