package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardBehindBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fast-forward a branch that is behind its upstream branch, without checking it out",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.NewBranch("feature")
		shell.EmptyCommit("two")
		shell.EmptyCommit("three")

		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("feature", "origin/feature")

		// remove the two commits so that there is something to fast-forward to
		shell.HardReset("HEAD~2")
		shell.Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("feature ↓2"),
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
				Contains("three"),
				Contains("two"),
				Contains("one"),
			)
	},
})
