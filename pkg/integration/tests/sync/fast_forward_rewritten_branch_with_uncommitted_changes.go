package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardRewrittenBranchWithUncommittedChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try to fast-forward a branch with a rewritten upstream branch while its worktree has uncommitted changes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content")
		createBranchRewrittenOnTheRemote(shell)

		shell.Checkout("feature")
		shell.UpdateFile("file", "changed")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("feature ↓2↑2").IsSelected(),
				Contains("master"),
			).
			Press(keys.Branches.FastForward)

		t.ExpectPopup().Alert().Title(Equals("Error")).
			Content(Contains("Cannot fast-forward 'feature' because the worktree")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("feature ↓2↑2").IsSelected(),
				Contains("master"),
			)

		t.Views().Files().
			Lines(
				Contains("file"),
			)
	},
})
