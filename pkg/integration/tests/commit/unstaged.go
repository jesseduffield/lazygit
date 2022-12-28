package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// TODO: find out why we can't use .SelectedLine() on the staging/stagingSecondary views.

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
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("myfile")).
			PressEnter()

		t.Views().Staging().
			IsFocused().
			Tap(func() {
				t.Views().StagingSecondary().Content(DoesNotContain("+myfile content"))
			}).
			// stage the first line
			PressPrimaryAction().
			Tap(func() {
				t.Views().Staging().Content(DoesNotContain("+myfile content"))
				t.Views().StagingSecondary().Content(Contains("+myfile content"))
			}).
			Press(keys.Files.CommitChanges)

		commitMessage := "my commit message"
		t.ExpectPopup().CommitMessagePanel().Type(commitMessage).Confirm()

		t.Views().Commits().
			Lines(
				Contains(commitMessage),
			)

		t.Views().Staging().IsFocused()

		// TODO: assert that the staging panel has been refreshed (it currently does not get correctly refreshed)
	},
})
