package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardOldFileChange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding a single file from an old commit (does rebase in background to remove the file but retain the other one)",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file0", "file0")
		shell.Commit("first commit")

		shell.CreateFileAndAdd("file1", "file2")
		shell.CreateFileAndAdd("fileToRemove", "fileToRemove")
		shell.Commit("commit to change")

		shell.CreateFileAndAdd("file3", "file3")
		shell.Commit("third commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("third commit").IsSelected(),
				Contains("commit to change"),
				Contains("first commit"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
				Contains("fileToRemove"),
			).
			SelectNextItem().
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard file changes")).
			Content(Contains("Are you sure you want to discard this commit's changes to this file?")).
			Confirm()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			)

		t.FileSystem().PathNotPresent("fileToRemove")
	},
})
