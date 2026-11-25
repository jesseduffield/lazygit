package interactive_rebase

import (
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

func handleConflictsFromSwap(t *TestDriver, expectedCommand string, selectConflict bool) {
	t.Common().AcknowledgeConflicts()

	// If the conflict comes from directly moving a commit, we want to keep the moved commit
	// selected, so selectConflict is false. In other cases (e.g. a conflict after "continue
	// rebase") we want to select the conflict commit.
	commitTwoMatcher := Contains("pick").Contains("commit two")
	conflictMatcher := Contains(expectedCommand).Contains("<-- CONFLICT --- commit three")
	if selectConflict {
		conflictMatcher.IsSelected()
	} else {
		commitTwoMatcher.IsSelected()
	}

	t.Views().Commits().
		Lines(
			Contains("─── Pending rebase todos"),
			commitTwoMatcher,
			conflictMatcher,
			Contains("─── Commits"),
			Contains("commit one"),
		)

	t.Views().Files().
		IsFocused().
		Lines(
			Contains("UU myfile"),
		).
		PressEnter()

	t.Views().MergeConflicts().
		IsFocused().
		TopLines(
			Contains("<<<<<<< HEAD"),
			Contains("one"),
			Contains("======="),
			Contains("three"),
			Contains(">>>>>>>"),
		).
		SelectNextItem().
		PressPrimaryAction() // pick "three"

	t.Common().ContinueOnConflictsResolved("rebase")

	t.Common().AcknowledgeConflicts()

	t.Views().Files().
		IsFocused().
		Lines(
			Contains("UU myfile"),
		).
		PressEnter()

	t.Views().MergeConflicts().
		IsFocused().
		TopLines(
			Contains("<<<<<<< HEAD"),
			Contains("three"),
			Contains("======="),
			Contains("two"),
			Contains(">>>>>>>"),
		).
		SelectNextItem().
		PressPrimaryAction() // pick "two"

	t.Common().ContinueOnConflictsResolved("rebase")

	t.Views().Commits().
		Focus().
		Lines(
			Contains("commit two").IsSelected(),
			Contains("commit three"),
			Contains("commit one"),
		).
		Tap(func() {
			t.Views().Main().Content(Contains("-three").Contains("+two"))
		}).
		SelectNextItem().
		Tap(func() {
			t.Views().Main().Content(Contains("-one").Contains("+three"))
		})
}

// Two branches forked from master whose commits interleave when sorted by
// date, so that with the whole graph shown a dev-a commit sits right below
// dev-b's first commit:
//
//	dev-a-2, dev-b-2, dev-b-1, dev-a-1, base
func setupInterleavedSiblingBranches(shell *Shell) {
	shell.
		EmptyCommitWithDate("base", "2026-01-01T00:00:00Z").
		NewBranch("dev-a").
		EmptyCommitWithDate("dev-a-1", "2026-01-01T00:01:00Z").
		EmptyCommitWithDate("dev-a-2", "2026-01-01T00:04:00Z").
		NewBranchFrom("dev-b", "master").
		EmptyCommitWithDate("dev-b-1", "2026-01-01T00:02:00Z").
		EmptyCommitWithDate("dev-b-2", "2026-01-01T00:03:00Z")
}
