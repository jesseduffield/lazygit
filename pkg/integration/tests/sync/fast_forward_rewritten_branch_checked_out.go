package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardRewrittenBranchCheckedOut = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fast-forward the checked out branch after its upstream branch was rewritten",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		createBranchRewrittenOnTheRemote(shell)

		shell.Checkout("feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("three"),
				Contains("two"),
				Contains("one"),
			)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("feature ↓2↑2").IsSelected(),
				Contains("master"),
			).
			Press(keys.Branches.FastForward).
			Lines(
				Contains("feature ✓").IsSelected(),
				Contains("master"),
			)

		t.Views().Commits().
			Lines(
				Contains("three-rewritten"),
				Contains("two-rewritten"),
				Contains("one"),
			)
	},
})
