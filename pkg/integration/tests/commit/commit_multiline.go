package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitMultiline = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with a multi-line commit message",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile content")
	},
	Run: func(shell *Shell, t *TestDriver, keys config.KeybindingConfig) {
		t.Model().CommitCount(0)

		t.Views().Files().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Files.CommitChanges)

		t.ExpectCommitMessagePanel().Type("first line").AddNewline().AddNewline().Type("third line").Confirm()

		t.Model().CommitCount(1)
		t.Model().HeadCommitMessage(Equals("first line"))

		t.Views().Commits().Focus()
		t.Views().Main().Content(MatchesRegexp("first line\n\\s*\n\\s*third line"))
	},
})
