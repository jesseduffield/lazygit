package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StagedWithoutHooks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files, going in the staged files menu, unstaging a line then committing without pre-commit hooks",
	ExtraCmdArgs: []string{},
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

		// stage the file
		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("myfile")).
			PressPrimaryAction().
			PressEnter()

		// we start with both lines having been staged
		t.Views().StagingSecondary().Content(
			Contains("+myfile content").Contains("+with a second line"),
		)
		t.Views().Staging().Content(
			DoesNotContain("+myfile content").DoesNotContain("+with a second line"),
		)

		// unstage the selected line
		t.Views().StagingSecondary().
			IsFocused().
			PressPrimaryAction().
			Tap(func() {
				// the line should have been moved to the main view
				t.Views().Staging().Content(Contains("+myfile content").DoesNotContain("+with a second line"))
			}).
			Content(DoesNotContain("+myfile content").Contains("+with a second line")).
			Press(keys.Files.CommitChangesWithoutHook)

		commitMessage := ": my commit message"
		t.ExpectPopup().CommitMessagePanel().InitialText(Contains("WIP")).Type(commitMessage).Confirm()

		t.Views().Commits().
			Lines(
				Contains("WIP" + commitMessage),
			)

		t.Views().StagingSecondary().
			IsEmpty()

		t.Views().Staging().
			IsFocused().
			Content(Contains("+myfile content")).
			Content(DoesNotContain("+with a second line"))
	},
})
