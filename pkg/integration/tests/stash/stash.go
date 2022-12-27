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
		assert.Model().StashCount(0)
		assert.Model().WorkingTreeFileCount(1)

		input.Press(keys.Files.ViewStashOptions)

		input.Menu().Title(Equals("Stash options")).Select(MatchesRegexp("stash all changes$")).Confirm()

		input.Prompt().Title(Equals("Stash changes")).Type("my stashed file").Confirm()

		assert.Model().StashCount(1)
		assert.Model().WorkingTreeFileCount(0)
	},
})
