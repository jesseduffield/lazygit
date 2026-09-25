package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckoutBranchOfOtherWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checkout a branch at a commit that is checked out by another worktree, which offers to switch to that worktree",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.AddWorktree("master", "../linked-worktree", "linked")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
			).
			PressPrimaryAction()

		t.ExpectPopup().Menu().
			Title(Contains("Checkout branch or commit")).
			Select(Contains("Checkout branch 'linked'")).
			Confirm()

		t.ExpectPopup().
			/* EXPECTED:
			Confirmation().
			Title(Equals("Switch to worktree")).
			Content(Equals("This branch is checked out by worktree linked-worktree. Do you want to switch to that worktree?")).
			ACTUAL: */
			Alert().
			Title(Equals("Error")).
			Content(Contains("already")).
			Confirm()

		t.Views().
			/* EXPECTED:
			Commits().
			ACTUAL: */
			Branches().
			IsFocused()

		t.Views().Branches().
			Lines(
				/* EXPECTED:
				Contains("linked"),
				Contains("master (worktree repo)"),
				ACTUAL: */
				Contains("master"),
				Contains("linked (worktree linked-worktree)"),
			)
	},
})
