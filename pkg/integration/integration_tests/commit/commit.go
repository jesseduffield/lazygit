package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/integration/types"
)

var Commit = types.NewTest(types.NewTestArgs{
	Description:  "Staging a couple files and committing",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell types.Shell) {
		shell.CreateFile("myfile", "myfile content")
		shell.CreateFile("myfile2", "myfile2 content")
	},
	Run: func(shell types.Shell, input types.Input, assert types.Assert, keys config.KeybindingConfig) {
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
