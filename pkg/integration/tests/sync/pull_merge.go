package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PullMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull with a merge strategy",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content1")
		shell.Commit("one")
		shell.UpdateFileAndAdd("file", "content2")
		shell.Commit("two")
		shell.EmptyCommit("three")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.HardReset("HEAD^^")
		shell.EmptyCommit("four")

		shell.SetConfig("pull.rebase", "false")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("four"),
				Contains("one"),
			)

		t.Views().Status().Content(Contains("↓2 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Pull)

		t.Views().Status().Content(Contains("↑2 repo → master"))

		t.Views().Commits().
			Lines(
				Contains("Merge branch 'master' of ../origin"),
				Contains("three"),
				Contains("two"),
				Contains("four"),
				Contains("one"),
			)
	},
})
