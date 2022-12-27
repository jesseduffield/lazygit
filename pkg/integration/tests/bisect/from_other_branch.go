package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FromOtherBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Opening lazygit when bisect has been started from another branch. There's an issue where we don't reselect the current branch if we mark the current branch as bad so this test side-steps that problem",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("only commit on master"). // this'll ensure we have a master branch
			NewBranch("other").
			CreateNCommits(10).
			Checkout("master").
			RunCommand("git bisect start other~2 other~5")
	},
	SetupConfig: func(cfg *config.AppConfig) {},
	Run: func(
		shell *Shell,
		input *Input,
		assert *Assert,
		keys config.KeybindingConfig,
	) {
		assert.View("information").Content(Contains("bisecting"))

		assert.AtLeastOneCommit()

		input.SwitchToCommitsView()

		assert.CurrentView().TopLines(
			MatchesRegexp(`<-- bad.*commit 08`),
			MatchesRegexp(`<-- current.*commit 07`),
			MatchesRegexp(`\?.*commit 06`),
			MatchesRegexp(`<-- good.*commit 05`),
		)

		input.NextItem()

		input.Press(keys.Commits.ViewBisectOptions)
		input.Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`mark .* as good`)).Confirm()

		input.Alert().Title(Equals("Bisect complete")).Content(MatchesRegexp("(?s)commit 08.*Do you want to reset")).Confirm()

		assert.View("information").Content(DoesNotContain("bisecting"))

		// back in master branch which just had the one commit
		assert.CurrentView().Name("commits").Lines(
			Contains("only commit on master"),
		)
	},
})
