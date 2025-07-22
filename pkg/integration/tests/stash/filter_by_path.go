package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterByPath = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter the stash list by path",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFileAndAdd("file1", "content")
		shell.Stash("file1")
		shell.CreateDir("subdir")
		shell.CreateFileAndAdd("subdir/file2", "content")
		shell.Stash("subdir/file2")
		shell.CreateFileAndAdd("file1", "other content")
		shell.Stash("file1 again")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		filterBy := func(path string) {
			t.GlobalPress(keys.Universal.FilteringMenu)
			t.ExpectPopup().Menu().
				Title(Equals("Filtering")).
				Select(Contains("Enter path to filter by")).
				Confirm()

			t.ExpectPopup().Prompt().
				Title(Equals("Enter path:")).
				Type(path).
				Confirm()
		}

		t.Views().Stash().
			Lines(
				Contains("file1 again"),
				Contains("subdir/file2"),
				Contains("file1"),
			)

		filterBy("file1")

		t.Views().Stash().
			Lines(
				Contains("file1 again"),
				Contains("file1"),
			)

		t.GlobalPress(keys.Universal.Return)
		filterBy("subdir")

		t.Views().Stash().
			Lines(
				Contains("subdir/file2"),
			)
	},
})
