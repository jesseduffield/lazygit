package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddFromCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add a worktree via the commits view",
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
			Press(keys.Worktrees.ViewWorktreeOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Worktree")).
					Select(MatchesRegexp(`Create worktree from .*`).DoesNotContain("detached")).
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
			Lines(
				Contains("initial commit"),
			)

		// Confirm we're now in the branches view
		t.Views().Branches().
			IsFocused().
			Lines(
				Contains("newbranch").IsSelected(),
				Contains("mybranch (worktree)"),
			)
	},
})
