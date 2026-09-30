package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectThemeWithBackgroundOverrides = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Selecting a theme with colors for dark and light backgrounds applies those for a dark one, which lazygit assumes when the terminal doesn't tell, as in integration tests",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", `
gui:
  theme:
    branchColorPatterns:
      master: '#0000ff'
  darkTheme:
    branchColorPatterns:
      master: '#ff00ff'
  lightTheme:
    branchColorPatterns:
      master: '#00ff00'
`)
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("master"),
			).
			DoesNotContainColoredText("#ff00ff", "master")

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Select(Contains("pink")).
			Confirm()

		t.ExpectToast(Equals("Theme: pink"))

		t.Views().Branches().
			ContainsColoredText("#ff00ff", "master").
			DoesNotContainColoredText("#0000ff", "master").
			DoesNotContainColoredText("#00ff00", "master")
	},
})
