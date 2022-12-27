package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiffCommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "View the diff between two commits",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "first line\n")
		shell.Commit("first commit")
		shell.UpdateFileAndAdd("file1", "first line\nsecond line\n")
		shell.Commit("second commit")
		shell.UpdateFileAndAdd("file1", "first line\nsecond line\nthird line\n")
		shell.Commit("third commit")
	},
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Views().Commits().
			Focus().
			Lines(
				Contains("third commit"),
				Contains("second commit"),
				Contains("first commit"),
			).
			Press(keys.Universal.DiffingMenu)

		input.ExpectMenu().Title(Equals("Diffing")).Select(MatchesRegexp(`diff \w+`)).Confirm()

		input.Views().Information().Content(Contains("showing output for: git diff"))

		input.Views().Commits().
			SelectNextItem().
			SelectNextItem().
			SelectedLine(Contains("first commit"))

		input.Views().Main().Content(Contains("-second line\n-third line"))

		input.Views().Commits().
			Press(keys.Universal.DiffingMenu)

		input.ExpectMenu().Title(Equals("Diffing")).Select(Contains("reverse diff direction")).Confirm()

		input.Views().Main().Content(Contains("+second line\n+third line"))

		input.Views().Commits().
			PressEnter()

		input.Views().CommitFiles().
			IsFocused().
			SelectedLine(Contains("file1"))

		input.Views().Main().Content(Contains("+second line\n+third line"))
	},
})
