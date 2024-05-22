package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddCoAuthorRange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add co-author on a range of commits",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("fourth commit")
		shell.EmptyCommit("third commit")
		shell.EmptyCommit("second commit")
		shell.EmptyCommit("first commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("first commit").IsSelected(),
				Contains("second commit"),
				Contains("third commit"),
				Contains("fourth commit"),
			).
			SelectNextItem().
			Press(keys.Universal.ToggleRangeSelect).
			SelectNextItem().
			Lines(
				Contains("first commit"),
				Contains("second commit").IsSelected(),
				Contains("third commit").IsSelected(),
				Contains("fourth commit"),
			).
			Press(keys.Commits.ResetCommitAuthor).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Amend commit attribute")).
					Select(Contains("Add co-author")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Contains("Add co-author")).
					Type("John Smith <jsmith@gmail.com>").
					Confirm()
			}).
			// exit range selection mode
			PressEscape().
			SelectNextItem()

		t.Views().Main().Content(
			Contains("fourth commit").
				DoesNotContain("Co-authored-by: John Smith <jsmith@gmail.com>"),
		)

		t.Views().Commits().
			IsFocused().
			SelectPreviousItem().
			Lines(
				Contains("first commit"),
				Contains("second commit"),
				Contains("third commit").IsSelected(),
				Contains("fourth commit"),
			)

		t.Views().Main().ContainsLines(
			Equals("    third commit"),
			Equals("    "),
			Equals("    Co-authored-by: John Smith <jsmith@gmail.com>"),
		)

		t.Views().Commits().
			IsFocused().
			SelectPreviousItem().
			Lines(
				Contains("first commit"),
				Contains("second commit").IsSelected(),
				Contains("third commit"),
				Contains("fourth commit"),
			)

		t.Views().Main().ContainsLines(
			Equals("    second commit"),
			Equals("    "),
			Equals("    Co-authored-by: John Smith <jsmith@gmail.com>"),
		)

		t.Views().Commits().
			IsFocused().
			SelectPreviousItem().
			Lines(
				Contains("first commit").IsSelected(),
				Contains("second commit"),
				Contains("third commit"),
				Contains("fourth commit"),
			)

		t.Views().Main().Content(
			Contains("first commit").
				DoesNotContain("Co-authored-by: John Smith <jsmith@gmail.com>"),
		)
	},
})
