package stash

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
)

var StashAll = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stashing all changes (via the menu)",
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
			Press(keys.Files.ViewStashOptions)

		t.ExpectPopup().Menu().Title(Equals("Stash options")).Select(MatchesRegexp("Stash all changes$")).Confirm()

		t.ExpectPopup().Prompt().Title(Equals("Stash changes")).Type("my stashed file").Confirm()

		t.Views().Stash().
			Lines(
				Contains("my stashed file"),
			)

		t.Views().Files().
			IsEmpty()
	},
})
