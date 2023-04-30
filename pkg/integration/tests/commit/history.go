package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var History = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cycling through commit message history in the commit message panel",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.EmptyCommit("commit 2")
		shell.EmptyCommit("commit 3")

		shell.CreateFile("myfile", "myfile content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			PressPrimaryAction(). // stage file
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("")).
			Type("my commit message").
			SelectPreviousMessage().
			Content(Equals("commit 3")).
			SelectPreviousMessage().
			Content(Equals("commit 2")).
			SelectPreviousMessage().
			Content(Equals("initial commit")).
			SelectPreviousMessage().
			Content(Equals("initial commit")). // we hit the end
			SelectNextMessage().
			Content(Equals("commit 2")).
			SelectNextMessage().
			Content(Equals("commit 3")).
			SelectNextMessage().
			Content(Equals("my commit message")).
			SelectNextMessage().
			Content(Equals("my commit message")). // we hit the beginning
			Type(" with extra added").
			Confirm()

		t.Views().Commits().
			TopLines(
				Contains("my commit message with extra added").IsSelected(),
			)
	},
})
