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
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Model().StashCount(0)
		t.Model().WorkingTreeFileCount(2)

		t.Views().Files().
			Press(keys.Files.ViewStashOptions)

		t.ExpectMenu().Title(Equals("Stash options")).Select(Contains("stash all changes including untracked files")).Confirm()

		t.ExpectPrompt().Title(Equals("Stash changes")).Type("my stashed file").Confirm()

		t.Model().StashCount(1)
		t.Model().WorkingTreeFileCount(0)
	},
})
