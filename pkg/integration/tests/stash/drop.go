package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Drop = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drop a stash entry",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFile("file", "content")
		shell.GitAddAll()
		shell.Stash("stash one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsEmpty()

		t.Views().Stash().
			Focus().
			Lines(
				Contains("stash one").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Stash drop")).
					Content(Contains("Are you sure you want to drop this stash entry?")).
					Confirm()
			}).
			IsEmpty()

		t.Views().Files().IsEmpty()
	},
})
