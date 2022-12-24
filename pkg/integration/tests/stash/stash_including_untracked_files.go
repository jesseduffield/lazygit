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
		assert.StashCount(0)
		assert.WorkingTreeFileCount(2)

		input.Press(keys.Files.ViewStashOptions)

		input.Menu(Equals("Stash options"), Contains("stash all changes including untracked files"))

		input.Prompt(Equals("Stash changes"), "my stashed file")

		assert.StashCount(1)
		assert.WorkingTreeFileCount(0)
	},
})
