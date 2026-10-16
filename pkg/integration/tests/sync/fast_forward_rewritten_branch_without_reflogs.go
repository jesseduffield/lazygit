package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardRewrittenBranchWithoutReflogs = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try to fast-forward a branch with a rewritten upstream branch in a repo that keeps no reflogs",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// Without reflogs there is no record of the values the remote-tracking
		// branch had before, so lazygit can't tell whether the commits that
		// the branch is ahead by were ever on the remote branch
		shell.SetConfig("core.logAllRefUpdates", "false")

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
			Press(keys.Branches.FastForward)

		t.ExpectPopup().Alert().Title(Equals("Error")).
			Content(Contains("Cannot fast-forward 'feature' because it has commits")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("master"),
				Contains("feature ↓2↑2").IsSelected(),
			)
	},
})
