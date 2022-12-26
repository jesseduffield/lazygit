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
		assert.View("commits").Lines(
			Contains("current-branch commit"),
			Contains("root commit"),
		)

		input.SwitchToBranchesWindow()

		assert.CurrentView().Name("localBranches").Lines(
			Contains("current-branch"),
			Contains("other-branch"),
		)
		input.NextItem()

		input.Press(keys.Commits.ViewResetOptions)

		input.Menu(Contains("reset to other-branch"), Contains("hard reset"))

		// ensure that we've returned from the menu before continuing
		assert.CurrentView().Name("localBranches")

		// assert that we now have the expected commits in the commit panel
		input.SwitchToCommitsWindow()
		assert.CurrentView().Name("commits").Lines(
			Contains("other-branch commit"),
			Contains("root commit"),
		)
	},
})
