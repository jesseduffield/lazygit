package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RevertWithConflictSingleCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reverts a commit that conflicts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("myfile", "")
		shell.Commit("add empty file")
		shell.CreateFileAndAdd("myfile", "first line\n")
		shell.Commit("add first line")
		shell.UpdateFileAndAdd("myfile", "first line\nsecond line\n")
		shell.Commit("add second line")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI ◯ add second line").IsSelected(),
				Contains("CI ◯ add first line"),
				Contains("CI ◯ add empty file"),
			).
			SelectNextItem().
			Press(keys.Commits.RevertCommit).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Revert commit")).
					Content(MatchesRegexp(`Are you sure you want to revert \w+?`)).
					Confirm()
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					// The exact error message is different on different git versions,
					// but they all contain the word 'conflict' somewhere.
					Content(Contains("conflict")).
					Confirm()
			}).
			Lines(
				/* EXPECTED:
				Proper display of revert commits is not implemented yet; we'll do this in the next PR
				Contains("revert").Contains("CI <-- CONFLICT --- add first line"),
				Contains("CI ◯ add second line"),
				Contains("CI ◯ add first line"),
				Contains("CI ◯ add empty file"),
				ACTUAL: */
				Contains("CI ◯ <-- YOU ARE HERE --- add second line"),
				Contains("CI ◯ add first line"),
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
				Contains(`CI ◯ Revert "add first line"`),
				Contains("CI ◯ add second line"),
				Contains("CI ◯ add first line"),
				Contains("CI ◯ add empty file"),
			)
	},
})
