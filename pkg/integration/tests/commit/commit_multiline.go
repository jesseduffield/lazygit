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
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Model().CommitCount(0)

		input.Views().Files().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Files.CommitChanges)

		input.ExpectCommitMessagePanel().Type("first line").AddNewline().AddNewline().Type("third line").Confirm()

		input.Model().CommitCount(1)
		input.Model().HeadCommitMessage(Equals("first line"))

		input.Views().Commits().Focus()
		input.Views().Main().Content(MatchesRegexp("first line\n\\s*\n\\s*third line"))
	},
})
