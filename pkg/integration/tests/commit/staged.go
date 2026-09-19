package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Staged = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files, going in the staged files menu, unstaging a line then committing",
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

		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Contains("myfile"),
				Contains("myfile2"),
			).
			SelectNextItem().
			PressPrimaryAction(). // stage the file
			Press(keys.Universal.FocusMainView)

		t.Views().Secondary().
			IsFocused().
			Tap(func() {
				// we start with both lines having been staged
				t.Views().Secondary().Content(Contains("+myfile content"))
				t.Views().Secondary().Content(Contains("+with a second line"))
				t.Views().Main().Content(DoesNotContain("+myfile content"))
				t.Views().Main().Content(DoesNotContain("+with a second line"))
			}).
			// unstage the selected line
			PressPrimaryAction().
			Tap(func() {
				// the line should have been moved to the main view
				t.Views().Secondary().Content(DoesNotContain("+myfile content"))
				t.Views().Secondary().Content(Contains("+with a second line"))
				t.Views().Main().Content(Contains("+myfile content"))
				t.Views().Main().Content(DoesNotContain("+with a second line"))
			}).
			Press(keys.Files.CommitChanges)

		commitMessage := "my commit message"
		t.ExpectPopup().CommitMessagePanel().Type(commitMessage).Confirm()

		t.Views().Commits().
			Lines(
				Contains(commitMessage),
			)

		t.Views().Secondary().IsInvisible()

		t.Views().Main().
			IsFocused().
			Content(Contains("+myfile content")).
			Content(DoesNotContain("+with a second line"))
	},
})
