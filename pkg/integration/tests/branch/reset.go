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
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Views().Commits().Lines(
			Contains("current-branch commit"),
			Contains("root commit"),
		)

		input.Views().Branches().
			Focus().
			Lines(
				Contains("current-branch"),
				Contains("other-branch"),
			).
			SelectNextItem().
			Press(keys.Commits.ViewResetOptions)

		input.ExpectMenu().Title(Contains("reset to other-branch")).Select(Contains("hard reset")).Confirm()

		// assert that we now have the expected commits in the commit panel
		input.Views().Commits().
			Lines(
				Contains("other-branch commit"),
				Contains("root commit"),
			)
	},
})
