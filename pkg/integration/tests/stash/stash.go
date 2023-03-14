package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Stash = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stashing files directly (not going through the stash menu)",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFile("file", "content")
		shell.GitAddAll()
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().
			IsEmpty()

		t.Views().Files().
			Lines(
				Contains("file"),
			).
			Press(keys.Files.StashAllChanges)

		t.ExpectPopup().Prompt().Title(Equals("Stash changes")).Type("my stashed file").Confirm()

		t.Views().Stash().
			Lines(
				Contains("my stashed file"),
			)

		t.Views().Files().
			IsEmpty()
	},
})
