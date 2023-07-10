package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PullRebaseInteractiveConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull with an interactive rebase strategy, where a conflict occurs",
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
		shell.UpdateFileAndAdd("file", "content4")
		shell.Commit("four")
		shell.CreateFileAndAdd("file5", "content5")
		shell.Commit("five")

		shell.SetConfig("pull.rebase", "interactive")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("five"),
				Contains("four"),
				Contains("one"),
			)

		t.Views().Status().Content(Contains("↓2 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Pull)

		t.Common().AcknowledgeConflicts()

		t.Views().Commits().
			Lines(
				Contains("pick").Contains("five"),
				Contains("conflict").Contains("YOU ARE HERE").Contains("four"),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			)

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

		t.Views().Status().Content(Contains("↑2 repo → master"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("five").IsSelected(),
				Contains("four"),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			SelectNextItem()

		t.Views().Main().
			Content(
				Contains("-content2").
					Contains("+content4"),
			)
	},
})
