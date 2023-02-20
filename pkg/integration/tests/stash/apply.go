package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Apply = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Apply a stash entry",
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
			PressPrimaryAction().
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Stash apply")).
					Content(Contains("Are you sure you want to apply this stash entry?")).
					Confirm()
			}).
			Lines(
				Contains("stash one").IsSelected(),
			)

		t.Views().Files().
			Lines(
				Contains("file"),
			)
	},
})
