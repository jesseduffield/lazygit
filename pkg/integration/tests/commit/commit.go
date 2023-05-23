package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Commit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files and committing",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile content")
		shell.CreateFile("myfile2", "myfile2 content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("?? myfile").IsSelected(),
				Contains("?? myfile2"),
			).
			PressPrimaryAction(). // stage file
			Lines(
				Contains("A  myfile").IsSelected(),
				Contains("?? myfile2"),
			).
			SelectNextItem().
			PressPrimaryAction(). // stage other file
			Lines(
				Contains("A  myfile"),
				Contains("A  myfile2").IsSelected(),
			).
			Press(keys.Files.CommitChanges)

		commitMessage := "my commit message"

		t.ExpectPopup().CommitMessagePanel().Type(commitMessage).Confirm()

		t.Views().Files().
			IsEmpty()

		t.Views().Commits().
			Focus().
			Lines(
				Contains(commitMessage).IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("A myfile"),
				Contains("A myfile2"),
			)
	},
})
