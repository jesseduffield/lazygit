package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var HistoryComplex = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "More complex flow for cycling commit message history",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.EmptyCommit("commit 2")
		shell.EmptyCommit("commit 3")

		shell.CreateFileAndAdd("myfile", "myfile content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// We're going to start a new commit message,
		// then leave and try to reword a commit, then
		// come back to original message and confirm we haven't lost our message.
		// This shows that we're storing the preserved message for a new commit separately
		// to the message when cycling history.

		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("")).
			Type("my commit message").
			Cancel()

		t.Views().Commits().
			Focus().
			SelectedLine(Contains("commit 3")).
			Press(keys.Commits.RenameCommit)

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("commit 3")).
			SelectNextMessage().
			Content(Equals("")).
			Type("reworded message").
			SelectPreviousMessage().
			Content(Equals("commit 3")).
			SelectNextMessage().
			Content(Equals("reworded message")).
			Cancel()

		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("my commit message"))
	},
})
