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
		viewBisectOptions := func() {
			input.PressKeys(keys.Commits.ViewBisectOptions)
			assert.InMenu()
		}
		markCommitAsBad := func() {
			viewBisectOptions()
			assert.MatchSelectedLine(Contains("bad"))

			input.Confirm()
		}

		markCommitAsGood := func() {
			viewBisectOptions()
			assert.MatchSelectedLine(Contains("bad"))
			input.NextItem()
			assert.MatchSelectedLine(Contains("good"))

			input.Confirm()
		}

		assert.AtLeastOneCommit()

		input.SwitchToCommitsWindow()

		assert.MatchSelectedLine(Contains("commit 10"))

		input.NavigateToListItemContainingText("commit 09")

		markCommitAsBad()

		assert.MatchViewContent("information", Contains("bisecting"))

		assert.CurrentViewName("commits")
		assert.MatchSelectedLine(Contains("<-- bad"))

		input.NavigateToListItemContainingText("commit 02")

		markCommitAsGood()

		// lazygit will land us in the comit between our good and bad commits.
		assert.CurrentViewName("commits")
		assert.MatchSelectedLine(Contains("commit 05"))
		assert.MatchSelectedLine(Contains("<-- current"))

		markCommitAsBad()

		assert.CurrentViewName("commits")
		assert.MatchSelectedLine(Contains("commit 04"))
		assert.MatchSelectedLine(Contains("<-- current"))

		markCommitAsGood()

		assert.InAlert()
		assert.MatchCurrentViewContent(Contains("Bisect complete!"))
		// commit 5 is the culprit because we marked 4 as good and 5 as bad.
		assert.MatchCurrentViewContent(Contains("commit 05"))
		assert.MatchCurrentViewContent(Contains("Do you want to reset"))
		input.Confirm()

		assert.CurrentViewName("commits")
		assert.MatchCurrentViewContent(Contains("commit 04"))
		assert.MatchViewContent("information", NotContains("bisecting"))
	},
})
