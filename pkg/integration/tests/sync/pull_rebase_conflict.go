package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PullRebaseConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull with a rebase strategy, where a conflict occurs",
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
		shell.UpdateFileAndAdd("file", "content4")
		shell.Commit("four")

		shell.SetConfig("pull.rebase", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("four"),
				Contains("one"),
			)

		t.Views().Status().Content(Equals("↓2↑1 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Pull)

		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU").Contains("file"),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			TopLines(
				Contains("<<<<<<< HEAD"),
				Contains("content2"),
				Contains("======="),
				Contains("content4"),
				Contains(">>>>>>>"),
			).
			SelectNextItem().
			PressPrimaryAction() // choose 'content4'

		t.Common().ContinueOnConflictsResolved()

		t.Views().Status().Content(Equals("↑1 repo → master"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			)

		t.Views().Main().
			Content(
				Contains("-content2").
					Contains("+content4"),
			)
	},
})
