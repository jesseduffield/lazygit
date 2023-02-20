package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Pop = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pop a stash entry",
	ExtraCmdArgs: "",
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
			Press(keys.Stash.PopStash).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Stash pop")).
					Content(Contains("Are you sure you want to pop this stash entry?")).
					Confirm()
			}).
			IsEmpty()

		t.Views().Files().
			Lines(
				Contains("file"),
			)
	},
})
