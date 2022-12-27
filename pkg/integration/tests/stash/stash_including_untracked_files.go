package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StashIncludingUntrackedFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stashing all files including untracked ones",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFile("file_1", "content")
		shell.CreateFile("file_2", "content")
		shell.GitAdd("file_1")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.Model().StashCount(0)
		assert.Model().WorkingTreeFileCount(2)

		input.Press(keys.Files.ViewStashOptions)

		input.Menu().Title(Equals("Stash options")).Select(Contains("stash all changes including untracked files")).Confirm()

		input.Prompt().Title(Equals("Stash changes")).Type("my stashed file").Confirm()

		assert.Model().StashCount(1)
		assert.Model().WorkingTreeFileCount(0)
	},
})
