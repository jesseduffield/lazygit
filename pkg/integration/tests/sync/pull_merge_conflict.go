package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PullMergeConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull with a merge strategy, where a conflict occurs",
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
				Contains("content4"),
				Contains("======="),
				Contains("content2"),
				Contains(">>>>>>>"),
			).
			PressPrimaryAction() // choose 'content4'

		t.Common().ContinueOnConflictsResolved()

		t.Views().Status().Content(Contains("↑2 repo → master"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("Merge branch 'master' of ../origin").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("four"),
				Contains("one"),
			)

		t.Views().Main().
			Content(
				Contains("- content4").
					Contains(" -content2").
					Contains("++content4"),
			)
	},
})
