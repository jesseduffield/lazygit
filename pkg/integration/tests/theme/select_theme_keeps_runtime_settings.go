package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectThemeKeepsRuntimeSettings = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Selecting a theme keeps settings that were changed at runtime, such as the branch sort order",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#ff00ff'\n")
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("commit").
			NewBranch("first").
			EmptyCommitWithDate("commit", "2023-04-07 10:00:00").
			NewBranch("second").
			EmptyCommitWithDate("commit", "2023-04-07 12:00:00").
			NewBranch("third").
			EmptyCommitWithDate("commit", "2023-04-07 11:00:00").
			Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Branches.SortOrder)

		t.ExpectPopup().Menu().
			Title(Equals("Sort order")).
			Select(Contains("Recency")).
			Confirm()

		t.Views().Branches().
			IsFocused().
			Lines(
				Contains("master").IsSelected(),
				Contains("third"),
				Contains("second"),
				Contains("first"),
			).
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Select(Contains("pink")).
			Confirm()

		t.ExpectToast(Equals("Theme: pink"))

		t.Views().Branches().
			IsFocused().
			ContainsColoredText("#ff00ff", "master").
			Lines(
				Contains("master").IsSelected(),
				Contains("third"),
				Contains("second"),
				Contains("first"),
			).
			Press(keys.Branches.SortOrder)

		t.ExpectPopup().Menu().
			Title(Equals("Sort order")).
			ContainsLines(
				Contains("r (•) Recency"),
				Contains("a ( ) Alphabetical"),
				Contains("d ( ) Date"),
			).
			Cancel()
	},
})
