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
		assert *Assert,
		keys config.KeybindingConfig,
	) {
		markCommitAsBad := func() {
			input.Press(keys.Commits.ViewBisectOptions)
			input.Menu(Equals("Bisect"), MatchesRegexp(`mark .* as bad`))
		}

		markCommitAsGood := func() {
			input.Press(keys.Commits.ViewBisectOptions)
			input.Menu(Equals("Bisect"), MatchesRegexp(`mark .* as good`))
		}

		assert.AtLeastOneCommit()

		input.SwitchToCommitsView()

		assert.CurrentView().SelectedLine(Contains("commit 10"))

		input.NavigateToListItem(Contains("commit 09"))

		markCommitAsBad()

		assert.View("information").Content(Contains("bisecting"))

		assert.CurrentView().Name("commits")
		assert.CurrentView().SelectedLine(Contains("<-- bad"))

		input.NavigateToListItem(Contains("commit 02"))

		markCommitAsGood()

		// lazygit will land us in the commit between our good and bad commits.
		assert.CurrentView().
			Name("commits").
			SelectedLine(Contains("commit 05")).
			SelectedLine(Contains("<-- current"))

		markCommitAsBad()

		assert.CurrentView().
			Name("commits").
			SelectedLine(Contains("commit 04")).
			SelectedLine(Contains("<-- current"))

		markCommitAsGood()

		// commit 5 is the culprit because we marked 4 as good and 5 as bad.
		input.Alert(Equals("Bisect complete"), MatchesRegexp("(?s)commit 05.*Do you want to reset"))

		assert.CurrentView().Name("commits").Content(Contains("commit 04"))
		assert.View("information").Content(NotContains("bisecting"))
	},
})
