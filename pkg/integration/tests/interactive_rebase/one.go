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
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToCommitsWindow()
		assert.CurrentView().Name("commits").Lines(
			Contains("commit 05"),
			Contains("commit 04"),
			Contains("commit 03"),
			Contains("commit 02"),
			Contains("commit 01"),
		)

		input.NavigateToListItem(Contains("commit 02"))
		input.Press(keys.Universal.Edit)

		assert.CurrentView().Lines(
			MatchesRegexp("pick.*commit 05"),
			MatchesRegexp("pick.*commit 04"),
			MatchesRegexp("pick.*commit 03"),
			MatchesRegexp("YOU ARE HERE.*commit 02"),
			Contains("commit 01"),
		)

		input.PreviousItem()
		input.Press(keys.Commits.MarkCommitAsFixup)
		assert.CurrentView().Lines(
			MatchesRegexp("pick.*commit 05"),
			MatchesRegexp("pick.*commit 04"),
			MatchesRegexp("fixup.*commit 03"),
			MatchesRegexp("YOU ARE HERE.*commit 02"),
			Contains("commit 01"),
		)

		input.PreviousItem()
		input.Press(keys.Universal.Remove)
		assert.CurrentView().Lines(
			MatchesRegexp("pick.*commit 05"),
			MatchesRegexp("drop.*commit 04"),
			MatchesRegexp("fixup.*commit 03"),
			MatchesRegexp("YOU ARE HERE.*commit 02"),
			Contains("commit 01"),
		)

		input.PreviousItem()
		input.Press(keys.Commits.SquashDown)

		assert.CurrentView().Lines(
			MatchesRegexp("squash.*commit 05"),
			MatchesRegexp("drop.*commit 04"),
			MatchesRegexp("fixup.*commit 03"),
			MatchesRegexp("YOU ARE HERE.*commit 02"),
			Contains("commit 01"),
		)

		input.ContinueRebase()

		assert.CurrentView().Lines(
			Contains("commit 02"),
			Contains("commit 01"),
		)
	},
})
