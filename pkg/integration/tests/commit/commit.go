package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Commit = components.NewIntegrationTest(components.NewIntegrationTestArgs{
	Description:  "Staging a couple files and committing",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *components.Shell) {
		shell.CreateFile("myfile", "myfile content")
		shell.CreateFile("myfile2", "myfile2 content")
	},
	Run: func(shell *components.Shell, input *components.Input, assert *components.Assert, keys config.KeybindingConfig) {
		assert.CommitCount(0)

		input.Select()
		input.NextItem()
		input.Select()
		input.PressKeys(keys.Files.CommitChanges)

		commitMessage := "my commit message"
		input.Type(commitMessage)
		input.Confirm()

		assert.CommitCount(1)
		assert.HeadCommitMessage(commitMessage)
	},
})
