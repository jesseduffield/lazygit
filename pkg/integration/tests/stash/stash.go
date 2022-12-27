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
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Model().StashCount(0)
		input.Model().WorkingTreeFileCount(1)

		input.Views().Files().
			Press(keys.Files.ViewStashOptions)

		input.ExpectMenu().Title(Equals("Stash options")).Select(MatchesRegexp("stash all changes$")).Confirm()

		input.ExpectPrompt().Title(Equals("Stash changes")).Type("my stashed file").Confirm()

		input.Model().StashCount(1)
		input.Model().WorkingTreeFileCount(0)
	},
})
