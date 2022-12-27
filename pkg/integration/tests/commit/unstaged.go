package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// TODO: find out why we can't use input.SelectedLine() on the staging/stagingSecondary views.

var Unstaged = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files, going in the unstaged files menu, staging a line and committing",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFile("myfile", "myfile content\nwith a second line").
			CreateFile("myfile2", "myfile2 content")
	},
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Model().CommitCount(0)

		input.Views().Files().
			IsFocused().
			SelectedLine(Contains("myfile")).
			PressEnter()

		input.Views().Staging().
			IsFocused().
			Tap(func() {
				input.Views().StagingSecondary().Content(DoesNotContain("+myfile content"))
			}).
			// stage the first line
			PressPrimaryAction().
			Tap(func() {
				input.Views().Staging().Content(DoesNotContain("+myfile content"))
				input.Views().StagingSecondary().Content(Contains("+myfile content"))
			}).
			Press(keys.Files.CommitChanges)

		commitMessage := "my commit message"
		input.ExpectCommitMessagePanel().Type(commitMessage).Confirm()

		input.Model().CommitCount(1)
		input.Model().HeadCommitMessage(Equals(commitMessage))
		input.Views().Staging().IsFocused()

		// TODO: assert that the staging panel has been refreshed (it currently does not get correctly refreshed)
	},
})
