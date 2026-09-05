package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SignOffRange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Sign off a range of commits",
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
			SelectedLines(
				Contains("second commit"),
				Contains("third commit"),
			).
			Press(keys.Commits.ResetCommitAuthor).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Amend commit attribute")).
					Select(Contains("Sign off")).
					Confirm()
			}).
			PressEscape().
			SelectNextItem().
			SelectedLine(Contains("fourth commit"))

		t.Views().Main().Content(
			Contains("fourth commit").
				DoesNotContain("Signed-off-by: CI <CI@example.com>"),
		)

		t.Views().Commits().
			IsFocused().
			SelectPreviousItem().
			SelectedLine(Contains("third commit"))

		t.Views().Main().ContainsLines(
			Equals("    third commit"),
			Equals("    "),
			Equals("    Signed-off-by: CI <CI@example.com>"),
		)

		t.Views().Commits().
			IsFocused().
			SelectPreviousItem().
			SelectedLine(Contains("second commit"))

		t.Views().Main().ContainsLines(
			Equals("    second commit"),
			Equals("    "),
			Equals("    Signed-off-by: CI <CI@example.com>"),
		)

		t.Views().Commits().
			IsFocused().
			SelectPreviousItem().
			SelectedLine(Contains("first commit"))

		t.Views().Main().Content(
			Contains("first commit").
				DoesNotContain("Signed-off-by: CI <CI@example.com>"),
		)
	},
})
