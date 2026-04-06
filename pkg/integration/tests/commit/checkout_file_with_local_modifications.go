package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckoutFileWithLocalModifications = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checkout a file from a commit that has local modifications",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("dir/file1.txt", "file1\n")
		shell.CreateFileAndAdd("dir/file2.txt", "file2\n")
		shell.Commit("one")
		shell.UpdateFile("dir/file1.txt", "file1\nfile1 change\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("one").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("  A file1.txt"),
				Equals("  A file2.txt"),
			).
			Press(keys.CommitFiles.CheckoutCommitFile)

		t.ExpectPopup().Alert().Title(Equals("Error")).
			Content(Contains("local modifications")).
			Confirm()
	},
})
