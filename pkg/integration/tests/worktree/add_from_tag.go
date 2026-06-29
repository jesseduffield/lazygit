package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddFromTag = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a detached worktree at a tag, entering a worktree name",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.CreateLightweightTag("v1.0", "HEAD")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Tags().
			Focus().
			Lines(
				Contains("v1.0").IsSelected(),
			).
			Press(keys.Worktrees.ViewWorktreeOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("New worktree")).
					Select(Contains("New detached worktree at 'v1.0'")).
					Confirm()

				// a tag has no good name to derive, so we're asked for one
				t.ExpectPopup().Prompt().
					Title(Equals("New worktree name")).
					Type("tag-worktree").
					Confirm()

				t.ExpectPopup().Menu().
					Title(Equals("Worktree location")).
					Confirm()
			})

		// we've switched into the new worktree, with a detached head
		t.Views().Branches().
			IsFocused().
			Lines(
				Contains("(no branch)").IsSelected(),
				Contains("mybranch (worktree repo)"),
			)

		t.Views().Status().
			Content(Contains("repo(tag-worktree)"))
	},
})
