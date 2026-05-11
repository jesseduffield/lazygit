package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Reset = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Hard reset to another branch",
	ExtraCmdArgs: []string{},
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
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Lines(
			Contains("current-branch commit"),
			Contains("root commit"),
		)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("current-branch").IsSelected(),
				Contains("other-branch"),
			).
			SelectNextItem().
			Press(keys.Commits.ViewResetOptions)

		t.ExpectPopup().Menu().
			Title(Contains("Reset to other-branch")).
			Select(Contains("Hard reset")).
			Confirm()

		// assert that we now have the expected commits in the commit panel
		t.Views().Commits().
			Lines(
				Contains("other-branch commit"),
				Contains("root commit"),
			)
	},
})
