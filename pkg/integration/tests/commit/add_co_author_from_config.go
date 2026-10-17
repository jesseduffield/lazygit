package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddCoAuthorFromConfig = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add a co-author configured in additionalAuthors",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.AdditionalAuthors = []string{
			"Jane Doe <jane@example.com>",
			"My AI Agent <agent@example.com>",
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.SetAuthor("John Doe", "john@doe.com")
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("initial commit").IsSelected(),
			).
			Press(keys.Commits.ResetCommitAuthor).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Amend commit attribute")).
					Select(Contains("Add co-author")).
					Confirm()

				// The configured authors show up as suggestions even though they've
				// never committed to this repo, alongside the authors from the
				// commit history, sorted together.
				t.ExpectPopup().Prompt().
					Title(Contains("Add co-author")).
					SuggestionLines(
						Equals("Jane Doe <jane@example.com>"),
						Equals("John Doe <john@doe.com>"),
						Equals("My AI Agent <agent@example.com>"),
					).
					ConfirmSuggestion(Equals("Jane Doe <jane@example.com>"))
			})

		t.Views().Main().ContainsLines(
			Equals("    initial commit"),
			Equals("    "),
			Equals("    Co-authored-by: Jane Doe <jane@example.com>"),
		)
	},
})
