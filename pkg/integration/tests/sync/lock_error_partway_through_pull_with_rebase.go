package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var LockErrorPartwayThroughPullWithRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull with a rebase strategy while another git process takes index.lock partway through the rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content1")
		shell.Commit("one")
		shell.UpdateFileAndAdd("file", "content2")
		shell.Commit("two")
		shell.CreateFileAndAdd("file3", "content3")
		shell.Commit("three")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.HardReset("HEAD^^")
		shell.CreateFileAndAdd("file4", "content4")
		shell.Commit("four")
		shell.CreateFileAndAdd("file5", "content5")
		shell.Commit("five")

		shell.SetConfig("pull.rebase", "true")

		// The rebase makes a commit when it picks "four"; the lock then stops
		// it when it tries to pick "five".
		shared.HoldIndexLockAfterNextCommit(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("five"),
				Contains("four"),
				Contains("one"),
			)

		t.Views().Status().Content(Equals("↓2↑2 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Pull)

		/* EXPECTED:
		t.Views().Status().Content(Equals("↑2 repo → master"))

		t.Views().Commits().
			Lines(
				Contains("five"),
				Contains("four"),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			)
		ACTUAL: */
		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("You are not currently on a branch.")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("pick").Contains("five"),
				Contains("─── Commits"),
				Contains("four"),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			)
	},
})
