package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AmendCommitWithConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Amends a staged file to a commit, causing a conflict there.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "1\n").Commit("one")
		shell.UpdateFileAndAdd("file", "1\n2\n").Commit("two")
		shell.UpdateFileAndAdd("file", "1\n2\n3\n").Commit("three")
		shell.UpdateFileAndAdd("file", "1\n2\n4\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			NavigateToLine(Contains("two")).
			Press(keys.Commits.AmendToCommit).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Amend commit")).
					Content(Contains("Are you sure you want to amend this commit with your staged files?")).
					Confirm()
				t.Common().AcknowledgeConflicts()
			}).
			Lines(
				Contains("pick").Contains("three"),
				Contains("conflict").Contains("<-- YOU ARE HERE --- fixup! two"),
				Contains("two"),
				Contains("one"),
			)

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file"),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			TopLines(
				Contains("1"),
				Contains("2"),
				Contains("<<<<<<< HEAD"),
				Contains("======="),
				Contains("4"),
				Contains(">>>>>>>"),
			).
			SelectNextItem().
			PressPrimaryAction() // pick "4"

		t.Common().ContinueOnConflictsResolved()

		t.Common().AcknowledgeConflicts()

		t.Views().Commits().
			Lines(
				Contains("<-- YOU ARE HERE --- three"),
				Contains("two"),
				Contains("one"),
			)
	},
})
