package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckoutFileFromCommitWithLocalChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checkout a file from a commit when there are local changes, showing a confirmation",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file.txt", "one\n")
		shell.Commit("one")
		shell.CreateFileAndAdd("file.txt", "two\n")
		shell.Commit("two")
		shell.CreateFileAndAdd("file.txt", "three\n")
		shell.Commit("three")
		// Create local uncommitted changes
		shell.UpdateFile("file.txt", "local changes\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			).
			NavigateToLine(Contains("two")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("M file.txt"),
			).
			Press(keys.CommitFiles.CheckoutCommitFile)

		// Should show confirmation dialog
		t.ExpectPopup().Confirmation().
			Title(Equals("Checkout file from commit")).
			Content(Contains("Are you sure you want to checkout this file? Your uncommitted changes will be lost.")).
			Confirm()

		// After confirmation, file should be checked out
		t.Views().Files().
			Lines(
				Equals("M  file.txt"),
			)

		t.FileSystem().FileContent("file.txt", Equals("two\n"))
	},
})
