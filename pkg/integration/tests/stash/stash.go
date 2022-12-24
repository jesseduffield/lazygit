package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Stash = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stashing files",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFile("file", "content")
		shell.GitAddAll()
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.StashCount(0)
		assert.WorkingTreeFileCount(1)

		input.Press(keys.Files.ViewStashOptions)

		input.Menu(Equals("Stash options"), MatchesRegexp("stash all changes$"))

		input.Prompt(Equals("Stash changes"), "my stashed file")

		assert.StashCount(1)
		assert.WorkingTreeFileCount(0)
	},
})
