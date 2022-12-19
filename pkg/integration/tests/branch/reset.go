package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Reset = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Hard reset to another branch",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("current-branch")
		shell.EmptyCommit("root commit")

		shell.NewBranch("other-branch")
		shell.EmptyCommit("other-branch commit")

		shell.Checkout("current-branch")
		shell.EmptyCommit("current-branch commit")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToBranchesWindow()
		assert.CurrentViewName("localBranches")

		assert.MatchSelectedLine(Contains("current-branch"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("other-branch"))

		input.PressKeys(keys.Commits.ViewResetOptions)
		assert.InMenu()
		assert.MatchCurrentViewTitle(Contains("reset to other-branch"))

		assert.MatchSelectedLine(Contains("soft reset"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("mixed reset"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("hard reset"))

		input.Confirm()

		// ensure that we've returned from the menu before continuing
		assert.CurrentViewName("localBranches")

		// assert that we now have the expected commits in the commit panel
		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")
		assert.CommitCount(2)
		assert.MatchSelectedLine(Contains("other-branch commit"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("root commit"))
	},
})
