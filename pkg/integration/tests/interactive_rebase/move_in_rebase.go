package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveInRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Via a single interactive rebase move a commit all the way up then back down then slightly back up again and apply the change",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(4)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 04").IsSelected(),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 01")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("commit 04"),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("YOU ARE HERE").Contains("commit 01").IsSelected(),
			).
			SelectPreviousItem().
			Press(keys.Commits.MoveUpCommit).
			Lines(
				Contains("commit 04"),
				Contains("commit 02").IsSelected(),
				Contains("commit 03"),
				Contains("YOU ARE HERE").Contains("commit 01"),
			).
			Press(keys.Commits.MoveUpCommit).
			Lines(
				Contains("commit 02").IsSelected(),
				Contains("commit 04"),
				Contains("commit 03"),
				Contains("YOU ARE HERE").Contains("commit 01"),
			).
			// assert we can't move past the top
			Press(keys.Commits.MoveUpCommit).
			Tap(func() {
				t.ExpectToast(Contains("Disabled: Cannot move any further"))
			}).
			Lines(
				Contains("commit 02").IsSelected(),
				Contains("commit 04"),
				Contains("commit 03"),
				Contains("YOU ARE HERE").Contains("commit 01"),
			).
			Press(keys.Commits.MoveDownCommit).
			Lines(
				Contains("commit 04"),
				Contains("commit 02").IsSelected(),
				Contains("commit 03"),
				Contains("YOU ARE HERE").Contains("commit 01"),
			).
			Press(keys.Commits.MoveDownCommit).
			Lines(
				Contains("commit 04"),
				Contains("commit 03"),
				Contains("commit 02").IsSelected(),
				Contains("YOU ARE HERE").Contains("commit 01"),
			).
			// assert we can't move past the bottom
			Press(keys.Commits.MoveDownCommit).
			Tap(func() {
				t.ExpectToast(Contains("Disabled: Cannot move any further"))
			}).
			Lines(
				Contains("commit 04"),
				Contains("commit 03"),
				Contains("commit 02").IsSelected(),
				Contains("YOU ARE HERE").Contains("commit 01"),
			).
			// move it back up one so that we land in a different order than we started with
			Press(keys.Commits.MoveUpCommit).
			Lines(
				Contains("commit 04"),
				Contains("commit 02").IsSelected(),
				Contains("commit 03"),
				Contains("YOU ARE HERE").Contains("commit 01"),
			).
			Tap(func() {
				t.Common().ContinueRebase()
			}).
			Lines(
				Contains("commit 04"),
				Contains("commit 02").IsSelected(),
				Contains("commit 03"),
				DoesNotContain("YOU ARE HERE").Contains("commit 01"),
			)
	},
})
