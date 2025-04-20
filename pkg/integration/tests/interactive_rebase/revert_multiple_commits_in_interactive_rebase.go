package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RevertMultipleCommitsInInteractiveRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reverts a range of commits, the first of which conflicts, in the middle of an interactive rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// TODO: use our revert UI once we support range-select for reverts
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "X",
				Context: "commits",
				Command: "git -c core.editor=: revert HEAD^ HEAD^^",
			},
		}
	},
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
			Press("X").
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					// The exact error message is different on different git versions,
					// but they all contain the word 'conflict' somewhere.
					Content(Contains("conflict")).
					Confirm()
			}).
			Lines(
				Contains("CI unrelated change 3"),
				Contains("CI unrelated change 2"),
				Contains("revert").Contains("CI unrelated change 1"),
				Contains("revert").Contains("CI <-- CONFLICT --- add first line"),
				Contains("CI ◯ add second line"),
				Contains("CI ◯ add first line"),
				Contains("CI ◯ unrelated change 1"),
				Contains("CI ◯ add empty file"),
			)

		t.Views().Options().Content(Contains("View revert options: m"))
		t.Views().Information().Content(Contains("Reverting (Reset)"))

		t.Views().Files().Focus().
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
				Contains("pick").Contains("CI unrelated change 3"),
				Contains("pick").Contains("CI unrelated change 2"),
				Contains(`CI ◯ <-- YOU ARE HERE --- Revert "unrelated change 1"`),
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
