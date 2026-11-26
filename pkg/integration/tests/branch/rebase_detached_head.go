package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebaseDetachedHead = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase a detached HEAD onto a branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("one").
			NewBranch("branch").
			EmptyCommit("branch commit").
			Checkout("master").
			EmptyCommit("two").
			Checkout("master^0")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("HEAD detached").IsSelected(),
				Contains("branch"),
				Contains("master"),
			).
			NavigateToLine(Contains("branch")).
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Rebase 'HEAD'")).
			Select(Contains("Simple rebase")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("two"),
				Contains("branch commit"),
				Contains("one"),
			)

		t.Views().Branches().
			Lines(
				Contains("HEAD detached"),
				Contains("branch"),
				Contains("master"),
			)
	},
})
