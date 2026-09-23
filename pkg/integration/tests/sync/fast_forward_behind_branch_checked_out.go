package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardBehindBranchCheckedOut = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fast-forward the checked out branch when it is behind its upstream branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content")
		shell.Commit("one")
		shell.EmptyCommit("two")
		shell.EmptyCommit("three")

		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")

		// remove the two commits so that there is something to fast-forward to
		shell.HardReset("HEAD~2")

		// a change that the fast-forward doesn't get in the way of
		shell.UpdateFile("file", "changed")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("one"),
			)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("master ↓2").IsSelected(),
			).
			Press(keys.Branches.FastForward).
			Lines(
				Contains("master ✓").IsSelected(),
			)

		t.Views().Commits().
			Lines(
				Contains("three"),
				Contains("two"),
				Contains("one"),
			)

		// Moving the branch forward doesn't touch the change
		t.Views().Files().
			Lines(
				Contains("file"),
			)
	},
})
