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
	Run: func(shell *Shell, t *TestDriver, keys config.KeybindingConfig) {
		t.Model().StashCount(0)
		t.Model().WorkingTreeFileCount(1)

		t.Views().Files().
			Press(keys.Files.ViewStashOptions)

		t.ExpectMenu().Title(Equals("Stash options")).Select(MatchesRegexp("stash all changes$")).Confirm()

		t.ExpectPrompt().Title(Equals("Stash changes")).Type("my stashed file").Confirm()

		t.Model().StashCount(1)
		t.Model().WorkingTreeFileCount(0)
	},
})
