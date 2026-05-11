package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var (
	initialFileContent = "first-line\nold-second-line\nthird-line\n"
	// We're indenting each line and modifying the second line
	updatedFileContent = " first-line\n new-second-line\n third-line\n"
)

var IgnoreWhitespace = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle whitespace in the diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("myfile", initialFileContent)
		shell.Commit("initial commit")
		shell.UpdateFile("myfile", updatedFileContent)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Main().ContainsLines(
			Contains(`-first-line`),
			Contains(`-old-second-line`),
			Contains(`-third-line`),
			Contains(`+ first-line`),
			Contains(`+ new-second-line`),
			Contains(`+ third-line`),
		)

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.ToggleWhitespaceInDiffView)

		// lines with only whitespace changes are ignored (first and third lines)
		t.Views().Main().ContainsLines(
			Contains(`  first-line`),
			Contains(`-old-second-line`),
			Contains(`+ new-second-line`),
			Contains(`  third-line`),
		)

		// when toggling again it goes back to showing whitespace
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.ToggleWhitespaceInDiffView)

		t.Views().Main().ContainsLines(
			Contains(`-first-line`),
			Contains(`-old-second-line`),
			Contains(`-third-line`),
			Contains(`+ first-line`),
			Contains(`+ new-second-line`),
			Contains(`+ third-line`),
		)
	},
})
