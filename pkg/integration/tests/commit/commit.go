package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Commit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files and committing",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile content")
		shell.CreateFile("myfile2", "myfile2 content")
	},
	Run: func(shell *Shell, t *TestDriver, keys config.KeybindingConfig) {
		t.Model().CommitCount(0)

		t.Views().Files().
			IsFocused().
			PressPrimaryAction(). // stage file
			SelectNextItem().
			PressPrimaryAction(). // stage other file
			Press(keys.Files.CommitChanges)

		commitMessage := "my commit message"

		t.ExpectCommitMessagePanel().Type(commitMessage).Confirm()

		t.Model().
			CommitCount(1).
			HeadCommitMessage(Equals(commitMessage))
	},
})
