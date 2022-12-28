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
		t.Views().Stash().
			IsEmpty()

		t.Views().Files().
			Lines(
				Contains("file_1"),
				Contains("file_2"),
			).
			Press(keys.Files.ViewStashOptions)

		t.ExpectPopup().Menu().Title(Equals("Stash options")).Select(Contains("stash all changes including untracked files")).Confirm()

		t.ExpectPopup().Prompt().Title(Equals("Stash changes")).Type("my stashed file").Confirm()

		t.Views().Stash().
			Lines(
				Contains("my stashed file"),
			)

		t.Views().Files().
			IsEmpty()
	},
})
