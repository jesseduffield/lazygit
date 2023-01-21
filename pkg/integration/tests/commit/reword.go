package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Reword = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files and committing",
	ExtraCmdArgs: "",
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
			PressPrimaryAction().
			Press(keys.Files.CommitChanges)

		commitMessage := "my commit message"

		t.ExpectPopup().CommitMessagePanel().Type(commitMessage).Confirm()
		t.Views().Commits().
			Lines(
				Contains(commitMessage),
			)

		t.Views().Files().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Files.CommitChanges)

		wipCommitMessage := "my commit message wip"

		t.ExpectPopup().CommitMessagePanel().Type(wipCommitMessage).Close()

		t.Views().Commits().Focus().
			Lines(
				Contains(commitMessage),
			).Press(keys.Commits.RenameCommit)

		t.ExpectPopup().CommitMessagePanel().
			SwitchToDescription().
			Type("some description").
			SwitchToSummary().
			Confirm()

		t.Views().Main().Content(MatchesRegexp("my commit message\n\\s*some description"))

		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().Confirm()
		t.Views().Commits().
			Lines(
				Contains(wipCommitMessage),
			)
	},
})
