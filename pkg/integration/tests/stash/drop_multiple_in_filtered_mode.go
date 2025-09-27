package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DropMultipleInFilteredMode = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drop multiple stash entries when filtering by path",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFileAndAdd("file1", "content1")
		shell.Stash("stash one")
		shell.CreateFileAndAdd("file2", "content2a")
		shell.Stash("stash two-a")
		shell.CreateFileAndAdd("file3", "content3")
		shell.Stash("stash three")
		shell.CreateFileAndAdd("file2", "content2b")
		shell.Stash("stash two-b")
		shell.CreateFileAndAdd("file4", "content4")
		shell.Stash("stash four")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().
			Lines(
				Contains("stash four"),
				Contains("stash two-b"),
				Contains("stash three"),
				Contains("stash two-a"),
				Contains("stash one"),
			)

		t.GlobalPress(keys.Universal.FilteringMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("Enter path to filter by")).
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("Enter path:")).
			Type("file2").
			Confirm()

		t.Views().Stash().
			Focus().
			Lines(
				Contains("stash two-b").IsSelected(),
				Contains("stash two-a"),
			).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Stash drop")).
					Content(Contains("Are you sure you want to drop the selected stash entry(ies)?")).
					Confirm()
			}).
			IsEmpty()

		t.GlobalPress(keys.Universal.Return) // cancel filtering mode
		t.Views().Stash().
			Lines(
				Contains("stash four"),
				Contains("stash three"),
				Contains("stash one"),
			)
	},
})
