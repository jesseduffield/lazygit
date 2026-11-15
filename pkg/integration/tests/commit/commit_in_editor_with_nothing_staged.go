package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitInEditorWithNothingStaged = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit using the git editor with nothing staged, confirming to stage all files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CreateFile("file1", "file1 content")
		shell.CreateFile("file2", "file2 content")

		shell.SetConfig("core.editor", "sh -c 'echo my commit message >.git/COMMIT_EDITMSG'")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("?? file1").IsSelected(),
				Equals("?? file2"),
			).
			Press(keys.Files.CommitChangesWithEditor).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("No files staged")).
					Content(Equals("You have not staged any files. Commit all files?")).
					Confirm()
			}).
			IsEmpty()

		t.Views().Commits().
			Lines(
				Contains("my commit message"),
				Contains("first commit"),
			)
	},
})
