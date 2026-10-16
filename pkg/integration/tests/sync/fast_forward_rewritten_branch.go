package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardRewrittenBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fast-forward a branch that has diverged from its upstream because the upstream branch was rewritten",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		createBranchRewrittenOnTheRemote(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("feature ↓2↑2"),
			).
			SelectNextItem().
			Press(keys.Branches.FastForward).
			Lines(
				Contains("master"),
				Contains("feature ✓").IsSelected(),
			).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("three-rewritten"),
				Contains("two-rewritten"),
				Contains("one"),
			)
	},
})
