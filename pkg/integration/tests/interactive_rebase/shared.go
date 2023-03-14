package interactive_rebase

import (
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

func handleConflictsFromSwap(t *TestDriver) {
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
			Contains("one"),
			Contains("======="),
			Contains("three"),
			Contains(">>>>>>>"),
		).
		SelectNextItem().
		PressPrimaryAction() // pick "three"

	t.Common().ContinueOnConflictsResolved()

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

	t.Common().ContinueOnConflictsResolved()

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
