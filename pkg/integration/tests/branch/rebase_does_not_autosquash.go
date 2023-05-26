package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebaseDoesNotAutosquash = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase a branch that has fixups onto another branch, and verify that the fixups are not squashed even if rebase.autoSquash is enabled globally.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.SetConfig("rebase.autoSquash", "true")

		shell.
			EmptyCommit("base").
			NewBranch("my-branch").
			Checkout("master").
			EmptyCommit("master commit").
			Checkout("my-branch").
			EmptyCommit("branch commit").
			EmptyCommit("fixup! branch commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("fixup! branch commit"),
				Contains("branch commit"),
				Contains("base"),
			)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("my-branch").IsSelected(),
				Contains("master"),
			).
			SelectNextItem().
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Rebase 'my-branch' onto 'master'")).
			Select(Contains("Simple rebase")).
			Confirm()

		t.Views().Commits().Lines(
			Contains("fixup! branch commit"),
			Contains("branch commit"),
			Contains("master commit"),
			Contains("base"),
		)
	},
})
