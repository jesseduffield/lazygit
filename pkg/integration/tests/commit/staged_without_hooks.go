package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StagedWithoutHooks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files, going in the staged files menu, unstaging a line then committing without pre-commit hooks",
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

		// stage the file
		input.Views().Files().
			IsFocused().
			SelectedLine(Contains("myfile")).
			PressPrimaryAction().
			PressEnter()

		// we start with both lines having been staged
		input.Views().StagingSecondary().Content(
			Contains("+myfile content").Contains("+with a second line"),
		)
		input.Views().Staging().Content(
			DoesNotContain("+myfile content").DoesNotContain("+with a second line"),
		)

		// unstage the selected line
		input.Views().StagingSecondary().
			IsFocused().
			PressPrimaryAction().
			Tap(func() {
				// the line should have been moved to the main view
				input.Views().Staging().Content(Contains("+myfile content").DoesNotContain("+with a second line"))
			}).
			Content(DoesNotContain("+myfile content").Contains("+with a second line")).
			Press(keys.Files.CommitChangesWithoutHook)

		commitMessage := ": my commit message"
		input.ExpectCommitMessagePanel().InitialText(Contains("WIP")).Type(commitMessage).Confirm()

		input.Model().CommitCount(1)
		input.Model().HeadCommitMessage(Equals("WIP" + commitMessage))
		input.Views().StagingSecondary().IsFocused()

		// TODO: assert that the staging panel has been refreshed (it currently does not get correctly refreshed)
	},
})
