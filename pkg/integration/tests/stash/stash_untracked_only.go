package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StashUntrackedOnly = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing 's' with only untracked files shows an error instead of silently doing nothing, since plain `git stash` ignores untracked files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFile("untracked", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().
			IsEmpty()

		t.Views().Files().
			Lines(
				Contains("untracked"),
			).
			Press(keys.Files.StashAllChanges)

		t.ExpectPopup().Alert().Title(Equals("Error")).Content(Equals("You have no files to stash")).Confirm()

		t.Views().Stash().
			IsEmpty()

		t.Views().Files().
			Lines(
				Contains("untracked"),
			)
	},
})
