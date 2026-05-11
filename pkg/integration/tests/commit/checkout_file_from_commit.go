package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckoutFileFromCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checkout a file from a commit",
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
		shell.CreateFileAndAdd("file.txt", "four\n")
		shell.Commit("four")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			NavigateToLine(Contains("three")).
			Tap(func() {
				t.Views().Main().ContainsLines(
					Contains("-two"),
					Contains("+three"),
				)
			}).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("M file.txt"),
			).
			Press(keys.CommitFiles.CheckoutCommitFile)

		t.Views().Files().
			Lines(
				Equals("M  file.txt"),
			)

		t.FileSystem().FileContent("file.txt", Equals("three\n"))
	},
})
