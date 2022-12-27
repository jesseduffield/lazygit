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
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.Model().CommitCount(0)

		input.PrimaryAction()
		input.NextItem()
		input.PrimaryAction()
		input.Press(keys.Files.CommitChanges)

		commitMessage := "my commit message"

		input.CommitMessagePanel().Type(commitMessage).Confirm()

		assert.Model().CommitCount(1)
		assert.Model().HeadCommitMessage(Equals(commitMessage))
	},
})
