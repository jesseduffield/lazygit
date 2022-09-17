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
		viewBisectOptions := func() {
			input.PressKeys(keys.Commits.ViewBisectOptions)
			assert.InMenu()
		}

		markCommitAsGood := func() {
			viewBisectOptions()
			assert.MatchSelectedLine(Contains("bad"))
			input.NextItem()
			assert.MatchSelectedLine(Contains("good"))

			input.Confirm()
		}

		assert.MatchViewContent("information", Contains("bisecting"))

		assert.AtLeastOneCommit()

		input.SwitchToCommitsWindow()

		assert.MatchSelectedLine(Contains("<-- bad"))
		assert.MatchSelectedLine(Contains("commit 08"))

		input.NextItem()
		assert.MatchSelectedLine(Contains("<-- current"))
		assert.MatchSelectedLine(Contains("commit 07"))

		markCommitAsGood()

		assert.InAlert()
		assert.MatchCurrentViewContent(Contains("Bisect complete!"))
		assert.MatchCurrentViewContent(Contains("commit 08"))
		assert.MatchCurrentViewContent(Contains("Do you want to reset"))
		input.Confirm()

		assert.MatchViewContent("information", NotContains("bisecting"))

		// back in master branch which just had the one commit
		assert.CurrentViewName("commits")
		assert.CommitCount(1)
		assert.MatchSelectedLine(Contains("only commit on master"))
	},
})
