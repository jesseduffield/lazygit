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
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToCommitsView()

		assert.CurrentView().Lines(
			Contains("third commit"),
			Contains("second commit"),
			Contains("first commit"),
		)

		input.Press(keys.Universal.DiffingMenu)
		input.Menu(Equals("Diffing"), MatchesRegexp(`diff \w+`))

		assert.NotInPopup()

		assert.View("information").Content(Contains("showing output for: git diff"))

		input.NextItem()
		input.NextItem()
		assert.CurrentView().SelectedLine(Contains("first commit"))

		assert.MainView().Content(Contains("-second line\n-third line"))

		input.Press(keys.Universal.DiffingMenu)
		input.Menu(Equals("Diffing"), Contains("reverse diff direction"))
		assert.NotInPopup()

		assert.MainView().Content(Contains("+second line\n+third line"))

		input.Enter()

		assert.CurrentView().Name("commitFiles")
		assert.CurrentView().SelectedLine(Contains("file1"))
		assert.MainView().Content(Contains("+second line\n+third line"))
	},
})
