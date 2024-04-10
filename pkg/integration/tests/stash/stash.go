package stash

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
)

var Stash = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stashing files directly (not going through the stash menu)",
	ExtraCmdArgs: []string{},
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
				MatchesRegexp(`\ds .* my stashed file`),
			)

		t.Views().Files().
			IsEmpty()
	},
})
