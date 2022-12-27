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
			input.Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`mark .* as bad`)).Confirm()
		}

		markCommitAsGood := func() {
			input.Press(keys.Commits.ViewBisectOptions)
			input.Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`mark .* as good`)).Confirm()
		}

		assert.Model().AtLeastOneCommit()

		input.SwitchToCommitsView()

		assert.Views().Current().SelectedLine(Contains("commit 10"))

		input.NavigateToListItem(Contains("commit 09"))

		markCommitAsBad()

		assert.Views().ByName("information").Content(Contains("bisecting"))

		assert.Views().Current().Name("commits").SelectedLine(Contains("<-- bad"))

		input.NavigateToListItem(Contains("commit 02"))

		markCommitAsGood()

		// lazygit will land us in the commit between our good and bad commits.
		assert.Views().Current().
			Name("commits").
			SelectedLine(Contains("commit 05").Contains("<-- current"))

		markCommitAsBad()

		assert.Views().Current().
			Name("commits").
			SelectedLine(Contains("commit 04").Contains("<-- current"))

		markCommitAsGood()

		// commit 5 is the culprit because we marked 4 as good and 5 as bad.
		input.Alert().Title(Equals("Bisect complete")).Content(MatchesRegexp("(?s)commit 05.*Do you want to reset")).Confirm()

		assert.Views().Current().Name("commits").Content(Contains("commit 04"))
		assert.Views().ByName("information").Content(DoesNotContain("bisecting"))
	},
})
