package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Unstaged = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files, going in the unstaged files menu, staging a line and committing",
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
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Tap(func() {
				t.Views().Secondary().Content(DoesNotContain("+myfile content"))
				t.Views().Main().SelectedLine(Equals("+myfile content"))
			}).
			// stage the first line
			PressPrimaryAction().
			Tap(func() {
				t.Views().Main().Content(DoesNotContain("+myfile content")).
					SelectedLine(Equals("+with a second line"))
				t.Views().Secondary().Content(Contains("+myfile content"))
			}).
			Press(keys.Files.CommitChanges)

		commitMessage := "my commit message"
		t.ExpectPopup().CommitMessagePanel().Type(commitMessage).Confirm()

		t.Views().Commits().
			Lines(
				Contains(commitMessage),
			)

		t.Views().Main().IsFocused()
	},
})
