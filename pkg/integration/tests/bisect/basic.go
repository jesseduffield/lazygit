package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Basic = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start a git bisect to find a bad commit",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(10)
	},
	SetupConfig: func(cfg *config.AppConfig) {},
	Run: func(
		shell *Shell,
		input *Input,
		keys config.KeybindingConfig,
	) {
		markCommitAsBad := func() {
			input.Views().Commits().
				Press(keys.Commits.ViewBisectOptions)

			input.ExpectMenu().Title(Equals("Bisect")).Select(MatchesRegexp(`mark .* as bad`)).Confirm()
		}

		markCommitAsGood := func() {
			input.Views().Commits().
				Press(keys.Commits.ViewBisectOptions)

			input.ExpectMenu().Title(Equals("Bisect")).Select(MatchesRegexp(`mark .* as good`)).Confirm()
		}

		input.Model().AtLeastOneCommit()

		input.Views().Commits().
			Focus().
			SelectedLine(Contains("commit 10")).
			NavigateToListItem(Contains("commit 09"))

		markCommitAsBad()

		input.Views().Information().Content(Contains("bisecting"))

		input.Views().Commits().
			IsFocused().
			SelectedLine(Contains("<-- bad")).
			NavigateToListItem(Contains("commit 02"))

		markCommitAsGood()

		// lazygit will land us in the commit between our good and bad commits.
		input.Views().Commits().IsFocused().
			SelectedLine(Contains("commit 05").Contains("<-- current"))

		markCommitAsBad()

		input.Views().Commits().IsFocused().
			SelectedLine(Contains("commit 04").Contains("<-- current"))

		markCommitAsGood()

		// commit 5 is the culprit because we marked 4 as good and 5 as bad.
		input.ExpectAlert().Title(Equals("Bisect complete")).Content(MatchesRegexp("(?s)commit 05.*Do you want to reset")).Confirm()

		input.Views().Commits().IsFocused().Content(Contains("commit 04"))
		input.Views().Information().Content(DoesNotContain("bisecting"))
	},
})
