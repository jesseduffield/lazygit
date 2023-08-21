package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitSwitchToEditor = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit, then switch from built-in commit message panel to editor",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("file1", "file1 content")
		shell.CreateFile("file2", "file2 content")

		// Set an editor that appends a line to the existing message. Since
		// git adds all this "# Please enter the commit message for your changes"
		// stuff, this will result in an extra blank line before the added line.
		shell.SetConfig("core.editor", "sh -c 'echo third line >>.git/COMMIT_EDITMSG'")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			PressPrimaryAction(). // stage one of the files
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("first line").
			SwitchToDescription().
			Type("second line").
			SwitchToSummary().
			SwitchToEditor()
		t.Views().Commits().
			Lines(
				Contains("first line"),
			)

		t.Views().Commits().Focus()
		t.Views().Main().Content(MatchesRegexp(`first line\n\s*\n\s*second line\n\s*\n\s*third line`))

		// Now check that the preserved commit message was cleared:
		t.Views().Files().
			Focus().
			PressPrimaryAction(). // stage the other file
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals(""))
	},
})
