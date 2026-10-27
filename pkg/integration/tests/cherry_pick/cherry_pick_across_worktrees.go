package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPickAcrossWorktrees = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy commits in one worktree, drop them there, and paste them in another worktree",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base").
			AddWorktree("master", "../linked-worktree", "newbranch").
			EmptyCommit("one").
			EmptyCommit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
				Contains("base"),
			).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Commits.CherryPickCopy).
			Tap(func() {
				t.Views().Information().Content(Contains("2 commits copied"))
			}).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Drop commit")).
					Content(Contains("Are you sure you want to drop the selected commit(s)?")).
					Confirm()
			}).
			Lines(
				Contains("base").IsSelected(),
			)

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("newbranch")).
			Press(keys.Universal.Select).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Switch to worktree")).
					Content(Equals("This branch is checked out by worktree linked-worktree. Do you want to switch to that worktree?")).
					Confirm()
			}).
			Lines(
				Contains("newbranch").IsSelected(),
				Contains("master (worktree repo)"),
			)

		t.Views().Information().Content(Contains("2 commits copied"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("base").IsSelected(),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Contains("Are you sure you want to cherry-pick the 2 copied commit(s) onto this branch?")).
					Confirm()
			}).
			Lines(
				Contains("two"),
				Contains("one"),
				Contains("base").IsSelected(),
			)
	},
})
