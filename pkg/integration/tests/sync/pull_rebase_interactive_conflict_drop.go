package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PullRebaseInteractiveConflictDrop = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull with an interactive rebase strategy, where a conflict occurs. Also drop a commit while rebasing",
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
		shell.CreateFileAndAdd("fil5", "content5")
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

		t.Views().Status().Content(Equals("↓2↑2 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Pull)

		t.Common().AcknowledgeConflicts()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("pick").Contains("five").IsSelected(),
				Contains("conflict").Contains("YOU ARE HERE").Contains("four"),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.Remove).
			Lines(
				Contains("drop").Contains("five").IsSelected(),
				Contains("conflict").Contains("YOU ARE HERE").Contains("four"),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			)

		t.Views().Files().
			Focus().
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
