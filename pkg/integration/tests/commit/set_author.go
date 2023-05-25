package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SetAuthor = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Set author on a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.SetConfig("user.email", "Bill@example.com")
		shell.SetConfig("user.name", "Bill Smith")

		shell.EmptyCommit("one")

		shell.SetConfig("user.email", "John@example.com")
		shell.SetConfig("user.name", "John Smith")

		shell.EmptyCommit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("JS").Contains("two").IsSelected(),
				Contains("BS").Contains("one"),
			).
			Press(keys.Commits.ResetCommitAuthor).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Amend commit attribute")).
					Select(Contains(" Set author")). // adding space at start to distinguish from 'reset author'
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Contains("Set author")).
					SuggestionLines(
						Contains("John Smith"),
						Contains("Bill Smith"),
					).
					ConfirmSuggestion(Contains("John Smith"))
			}).
			Lines(
				Contains("JS").Contains("two").IsSelected(),
				Contains("BS").Contains("one"),
			)
	},
})
