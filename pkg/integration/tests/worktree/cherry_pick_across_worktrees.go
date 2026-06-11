package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPickAcrossWorktrees = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy a commit in one worktree and paste it in another worktree of the same repo",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.EmptyCommit("base")
		// the linked worktree's branch stays at "base"
		shell.AddWorktree("mybranch", "../linked-worktree", "newbranch")
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("one").IsSelected(),
				Contains("base"),
			).
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("1 commit copied"))

		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("(main worktree)").IsSelected(),
				Contains("linked-worktree"),
			).
			NavigateToLine(Contains("linked-worktree")).
			Press(keys.Universal.Select)

		t.Views().Commits().
			Focus().
			Lines(
				Contains("base"),
			).
			Press(keys.Commits.PasteCommits)

		/* EXPECTED:
		t.ExpectPopup().Alert().
			Title(Equals("Cherry-pick")).
			Content(Contains("Are you sure you want to cherry-pick the 1 copied commit(s) onto this branch?")).
			Confirm()

		t.Views().Information().Content(DoesNotContain("commit copied"))

		t.Views().Commits().
			Lines(
				Contains("one"),
				Contains("base").IsSelected(),
			)
		ACTUAL: */
		t.ExpectToast(Equals("Disabled: No copied commits"))
	},
})
