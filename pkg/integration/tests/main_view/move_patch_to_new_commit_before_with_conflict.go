package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePatchToNewCommitBeforeWithConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch into a new commit before the original commit, where a commit hook changes the new commit so that replaying the original commit conflicts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "1st line\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "1st line\n2nd line\n")
		shell.CreateFileAndAdd("file2", "file2 content\n")
		shell.Commit("commit to move from")

		shell.CreateFile(".git/hooks/pre-commit", "#!/bin/sh\necho 'hook line' >> file1\ngit add file1\n")
		shell.MakeExecutable(".git/hooks/pre-commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit to move from").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Contains("file1"),
				Contains("file2"),
			).
			NavigateToLine(Contains("file2")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Common().SelectPatchOption(Contains("Move patch into new commit before the original commit"))

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("")).
			Type("new commit").Confirm()

		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU").Contains("file1"),
			)
	},
})
