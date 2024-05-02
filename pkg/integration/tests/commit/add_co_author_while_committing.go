package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddCoAuthorWhileCommitting = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add co-author while typing the commit message",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("file", "file content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			PressPrimaryAction(). // stage file
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("Subject").
			SwitchToDescription().
			Type("Here's my message.").
			AddCoAuthor("John Doe <john@doe.com>").
			Content(Equals("Here's my message.\n\nCo-authored-by: John Doe <john@doe.com>")).
			AddCoAuthor("Jane Smith <jane@smith.com>").
			// Second co-author doesn't add a blank line:
			Content(Equals("Here's my message.\n\nCo-authored-by: John Doe <john@doe.com>\nCo-authored-by: Jane Smith <jane@smith.com>")).
			SwitchToSummary().
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("Subject"),
			).
			Focus().
			Tap(func() {
				t.Views().Main().ContainsLines(
					Equals("    Subject"),
					Equals("    "),
					Equals("    Here's my message."),
					Equals("    "),
					Equals("    Co-authored-by: John Doe <john@doe.com>"),
					Equals("    Co-authored-by: Jane Smith <jane@smith.com>"),
				)
			})
	},
})
