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
			assert.SelectedLine(Contains("bad"))

			input.Confirm()
		}

		markCommitAsGood := func() {
			viewBisectOptions()
			assert.SelectedLine(Contains("bad"))
			input.NextItem()
			assert.SelectedLine(Contains("good"))

			input.Confirm()
		}

		assert.AtLeastOneCommit()

		input.SwitchToCommitsWindow()

		assert.SelectedLine(Contains("commit 10"))

		input.NavigateToListItemContainingText("commit 09")

		markCommitAsBad()

		assert.ViewContent("information", Contains("bisecting"))

		assert.CurrentViewName("commits")
		assert.SelectedLine(Contains("<-- bad"))

		input.NavigateToListItemContainingText("commit 02")

		markCommitAsGood()

		// lazygit will land us in the comit between our good and bad commits.
		assert.CurrentViewName("commits")
		assert.SelectedLine(Contains("commit 05"))
		assert.SelectedLine(Contains("<-- current"))

		markCommitAsBad()

		assert.CurrentViewName("commits")
		assert.SelectedLine(Contains("commit 04"))
		assert.SelectedLine(Contains("<-- current"))

		markCommitAsGood()

		assert.InAlert()
		assert.CurrentViewContent(Contains("Bisect complete!"))
		// commit 5 is the culprit because we marked 4 as good and 5 as bad.
		assert.CurrentViewContent(Contains("commit 05"))
		assert.CurrentViewContent(Contains("Do you want to reset"))
		input.Confirm()

		assert.CurrentViewName("commits")
		assert.CurrentViewContent(Contains("commit 04"))
		assert.ViewContent("information", NotContains("bisecting"))
	},
})
