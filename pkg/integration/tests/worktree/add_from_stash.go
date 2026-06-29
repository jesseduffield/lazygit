package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddFromStash = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a new branch and worktree from a stash entry",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.UpdateFile("README.md", "work in progress")
		shell.Stash("my stash")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().
			Focus().
			Lines(
				Contains("my stash").IsSelected(),
			).
			Press(keys.Worktrees.ViewWorktreeOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("New worktree")).
					Select(Contains("New branch and worktree from 'stash@{0}'")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New branch and worktree name")).
					Type("from-stash").
					Confirm()

				t.ExpectPopup().Menu().
					Title(Equals("Worktree location")).
					Confirm()
			})

		// we've switched into the new worktree, on the new branch
		t.Views().Branches().
			IsFocused().
			Lines(
				Contains("from-stash").IsSelected(),
				Contains("mybranch (worktree repo)"),
			)
	},
})
