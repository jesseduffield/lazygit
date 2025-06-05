package interactive_rebase

import (
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

func handleConflictsFromSwap(t *TestDriver, expectedCommand string) {
	t.Common().AcknowledgeConflicts()

	t.Views().Commits().
		Lines(
			Contains("--- Pending rebase todos ---"),
			Contains("pick").Contains("commit two"),
			Contains(expectedCommand).Contains("<-- CONFLICT --- commit three"),
			Contains("--- Commits ---"),
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
