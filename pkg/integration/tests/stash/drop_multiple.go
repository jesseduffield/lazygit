package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DropMultiple = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drop multiple stash entries",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFileAndAdd("file1", "content1")
		shell.Stash("stash one")
		shell.CreateFileAndAdd("file2", "content2")
		shell.Stash("stash two")
		shell.CreateFileAndAdd("file3", "content3")
		shell.Stash("stash three")
		shell.CreateFileAndAdd("file4", "content4")
		shell.Stash("stash four")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsEmpty()

		t.Views().Stash().
			Focus().
			SelectNextItem().
			Lines(
				Contains("stash four"),
				Contains("stash three").IsSelected(),
				Contains("stash two"),
				Contains("stash one"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Stash drop")).
					Content(Contains("Are you sure you want to drop the selected stash entry(ies)?")).
					Confirm()
			}).
			Lines(
				Contains("stash four"),
				Contains("stash one"),
			)

		t.Views().Files().IsEmpty()
	},
})
