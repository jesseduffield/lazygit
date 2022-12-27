package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var One = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Begins an interactive rebase, then fixups, drops, and squashes some commits",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(5) // these will appears at commit 05, 04, 04, down to 01
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 05"),
				Contains("commit 04"),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToListItem(Contains("commit 02")).
			Press(keys.Universal.Edit).
			Lines(
				MatchesRegexp("pick.*commit 05"),
				MatchesRegexp("pick.*commit 04"),
				MatchesRegexp("pick.*commit 03"),
				MatchesRegexp("YOU ARE HERE.*commit 02").IsSelected(),
				Contains("commit 01"),
			).
			SelectPreviousItem().
			Press(keys.Commits.MarkCommitAsFixup).
			Lines(
				MatchesRegexp("pick.*commit 05"),
				MatchesRegexp("pick.*commit 04"),
				MatchesRegexp("fixup.*commit 03").IsSelected(),
				MatchesRegexp("YOU ARE HERE.*commit 02"),
				Contains("commit 01"),
			).
			SelectPreviousItem().
			Press(keys.Universal.Remove).
			Lines(
				MatchesRegexp("pick.*commit 05"),
				MatchesRegexp("drop.*commit 04").IsSelected(),
				MatchesRegexp("fixup.*commit 03"),
				MatchesRegexp("YOU ARE HERE.*commit 02"),
				Contains("commit 01"),
			).
			SelectPreviousItem().
			Press(keys.Commits.SquashDown).
			Lines(
				MatchesRegexp("squash.*commit 05").IsSelected(),
				MatchesRegexp("drop.*commit 04"),
				MatchesRegexp("fixup.*commit 03"),
				MatchesRegexp("YOU ARE HERE.*commit 02"),
				Contains("commit 01"),
			).
			Tap(func() {
				t.Actions().ContinueRebase()
			}).
			Lines(
				Contains("commit 02"),
				Contains("commit 01"),
			)
	},
})
