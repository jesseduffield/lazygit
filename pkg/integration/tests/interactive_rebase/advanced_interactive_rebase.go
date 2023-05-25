package interactive_rebase

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

const (
	BASE_BRANCH = "base-branch"
	TOP_BRANCH  = "top-branch"
	BASE_COMMIT = "base-commit"
	TOP_COMMIT  = "top-commit"
)

var AdvancedInteractiveRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "It begins an interactive rebase and verifies to have the possibility of editing the commits of the branch before proceeding with the actual rebase",
	ExtraCmdArgs: []string{},
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch(BASE_BRANCH).
			EmptyCommit(BASE_COMMIT).
			NewBranch(TOP_BRANCH).
			EmptyCommit(TOP_COMMIT)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains(TOP_COMMIT),
				Contains(BASE_COMMIT),
			)

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains(BASE_BRANCH)).
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Menu().
			Title(Equals(fmt.Sprintf("Rebase '%s' onto '%s'", TOP_BRANCH, BASE_BRANCH))).
			Select(Contains("Interactive rebase")).
			Confirm()
		t.Views().Commits().
			IsFocused().
			Lines(
				Contains(TOP_COMMIT),
				Contains(BASE_COMMIT).Contains("YOU ARE HERE"),
			).
			NavigateToLine(Contains(TOP_COMMIT)).
			Press(keys.Universal.Edit).
			Lines(
				Contains(TOP_COMMIT).Contains("edit"),
				Contains(BASE_COMMIT).Contains("YOU ARE HERE"),
			).
			Tap(func() {
				t.Common().ContinueRebase()
			}).
			Lines(
				Contains(TOP_COMMIT).Contains("YOU ARE HERE"),
				Contains(BASE_COMMIT),
			)
	},
})
