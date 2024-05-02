package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OutsideRebaseRangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Do various things with range selection in the commits view when outside rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(10)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			TopLines(
				Contains("commit 10").IsSelected(),
			).
			Press(keys.Universal.RangeSelectDown).
			TopLines(
				Contains("commit 10").IsSelected(),
				Contains("commit 09").IsSelected(),
				Contains("commit 08"),
			).
			// Drop commits
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Drop commit")).
					Content(Contains("Are you sure you want to drop the selected commit(s)?")).
					Confirm()
			}).
			TopLines(
				Contains("commit 08").IsSelected(),
				Contains("commit 07"),
			).
			Press(keys.Universal.RangeSelectDown).
			TopLines(
				Contains("commit 08").IsSelected(),
				Contains("commit 07").IsSelected(),
				Contains("commit 06"),
			).
			// Squash commits
			Press(keys.Commits.SquashDown).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Squash")).
					Content(Contains("Are you sure you want to squash the selected commit(s) into the commit below?")).
					Confirm()
			}).
			TopLines(
				Contains("commit 06").IsSelected(),
				Contains("commit 05"),
				Contains("commit 04"),
			).
			// Verify commit messages are concatenated
			Tap(func() {
				t.Views().Main().
					ContainsLines(
						Contains("commit 06"),
						AnyString(),
						Contains("commit 07"),
						AnyString(),
						Contains("commit 08"),
					)
			}).
			// Fixup commits
			Press(keys.Universal.RangeSelectDown).
			TopLines(
				Contains("commit 06").IsSelected(),
				Contains("commit 05").IsSelected(),
				Contains("commit 04"),
			).
			Press(keys.Commits.MarkCommitAsFixup).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Fixup")).
					Content(Contains("Are you sure you want to 'fixup' the selected commit(s) into the commit below?")).
					Confirm()
			}).
			TopLines(
				Contains("commit 04").IsSelected(),
				Contains("commit 03"),
				Contains("commit 02"),
			).
			// Verify commit messages are dropped
			Tap(func() {
				t.Views().Main().
					Content(
						Contains("commit 04").
							DoesNotContain("commit 06").
							DoesNotContain("commit 05"),
					)
			}).
			Press(keys.Universal.RangeSelectDown).
			TopLines(
				Contains("commit 04").IsSelected(),
				Contains("commit 03").IsSelected(),
				Contains("commit 02"),
			).
			// Move commits
			Press(keys.Commits.MoveDownCommit).
			TopLines(
				Contains("commit 02"),
				Contains("commit 04").IsSelected(),
				Contains("commit 03").IsSelected(),
				Contains("commit 01"),
			).
			Press(keys.Commits.MoveDownCommit).
			TopLines(
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("commit 04").IsSelected(),
				Contains("commit 03").IsSelected(),
			).
			Press(keys.Commits.MoveDownCommit).
			TopLines(
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("commit 04").IsSelected(),
				Contains("commit 03").IsSelected(),
			).
			Tap(func() {
				t.ExpectToast(Contains("Disabled: Cannot move any further"))
			}).
			Press(keys.Commits.MoveUpCommit).
			TopLines(
				Contains("commit 02"),
				Contains("commit 04").IsSelected(),
				Contains("commit 03").IsSelected(),
				Contains("commit 01"),
			).
			Press(keys.Commits.MoveUpCommit).
			TopLines(
				Contains("commit 04").IsSelected(),
				Contains("commit 03").IsSelected(),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			Press(keys.Commits.MoveUpCommit).
			Tap(func() {
				t.ExpectToast(Contains("Disabled: Cannot move any further"))
			}).
			TopLines(
				Contains("commit 04").IsSelected(),
				Contains("commit 03").IsSelected(),
				Contains("commit 02"),
				Contains("commit 01"),
			)
	},
})
