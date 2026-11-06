package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StashNoTrackedChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing stash with only untracked files shows an error",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("untracked", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Lines(
				Contains("untracked"),
			).
			Press(keys.Files.StashAllChanges)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("You have no files to stash")).
			Confirm()

		t.Views().Files().
			Lines(
				Contains("untracked"),
			)
	},
})
