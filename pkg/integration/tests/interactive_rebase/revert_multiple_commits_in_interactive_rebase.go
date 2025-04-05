package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RevertMultipleCommitsInInteractiveRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reverts a range of commits, the first of which conflicts, in the middle of an interactive rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("myfile", "")
		shell.Commit("add empty file")
		shell.CreateFileAndAdd("otherfile", "")
		shell.Commit("unrelated change 1")
		shell.CreateFileAndAdd("myfile", "first line\n")
		shell.Commit("add first line")
		shell.UpdateFileAndAdd("myfile", "first line\nsecond line\n")
		shell.Commit("add second line")
		shell.EmptyCommit("unrelated change 2")
		shell.EmptyCommit("unrelated change 3")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI ◯ unrelated change 3").IsSelected(),
				Contains("CI ◯ unrelated change 2"),
				Contains("CI ◯ add second line"),
				Contains("CI ◯ add first line"),
				Contains("CI ◯ unrelated change 1"),
				Contains("CI ◯ add empty file"),
			).
			NavigateToLine(Contains("add second line")).
			Press(keys.Universal.Edit).
			SelectNextItem().
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Commits.RevertCommit).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Revert commit")).
					Content(Equals("Are you sure you want to revert the selected commits?")).
					Confirm()

				t.ExpectPopup().Menu().
					Title(Equals("Conflicts!")).
					Select(Contains("View conflicts")).
					Confirm()
			}).
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("CI unrelated change 3"),
				Contains("CI unrelated change 2"),
				Contains("--- Pending reverts ---"),
				Contains("revert").Contains("CI unrelated change 1"),
				Contains("revert").Contains("CI <-- CONFLICT --- add first line"),
				Contains("--- Commits ---"),
				Contains("CI ◯ add second line"),
				Contains("CI ◯ add first line"),
				Contains("CI ◯ unrelated change 1"),
				Contains("CI ◯ add empty file"),
			)

		t.Views().Options().Content(Contains("View revert options: m"))
		t.Views().Information().Content(Contains("Reverting (Reset)"))

		t.Views().Files().IsFocused().
			Lines(
				Contains("UU myfile").IsSelected(),
			).
			PressEnter()

		t.Views().MergeConflicts().IsFocused().
			SelectNextItem().
			PressPrimaryAction()

		t.ExpectPopup().Alert().
			Title(Equals("Continue")).
			Content(Contains("All merge conflicts resolved. Continue the revert?")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("pick").Contains("CI unrelated change 3"),
				Contains("pick").Contains("CI unrelated change 2"),
				Contains("--- Commits ---"),
				Contains(`CI ◯ Revert "unrelated change 1"`),
				Contains(`CI ◯ Revert "add first line"`),
				Contains("CI ◯ add second line"),
				Contains("CI ◯ add first line"),
				Contains("CI ◯ unrelated change 1"),
				Contains("CI ◯ add empty file"),
			)

		t.Views().Options().Content(Contains("View rebase options: m"))
		t.Views().Information().Content(Contains("Rebasing (Reset)"))

		t.Common().ContinueRebase()

		t.Views().Commits().
			Lines(
				Contains("CI ◯ unrelated change 3"),
				Contains("CI ◯ unrelated change 2"),
				Contains(`CI ◯ Revert "unrelated change 1"`),
				Contains(`CI ◯ Revert "add first line"`),
				Contains("CI ◯ add second line"),
				Contains("CI ◯ add first line"),
				Contains("CI ◯ unrelated change 1"),
				Contains("CI ◯ add empty file"),
			)
	},
})
