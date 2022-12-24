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
		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		assert.SelectedLine(Contains("third commit"))

		input.Press(keys.Universal.DiffingMenu)
		input.Menu(Equals("Diffing"), MatchesRegexp(`diff \w+`))

		assert.NotInPopup()

		assert.ViewContent("information", Contains("showing output for: git diff"))

		input.NextItem()
		input.NextItem()

		assert.SelectedLine(Contains("first commit"))

		assert.MainViewContent(Contains("-second line\n-third line"))

		input.Press(keys.Universal.DiffingMenu)
		input.Menu(Equals("Diffing"), Contains("reverse diff direction"))
		assert.NotInPopup()

		assert.MainViewContent(Contains("+second line\n+third line"))

		input.Enter()

		assert.CurrentViewName("commitFiles")
		assert.SelectedLine(Contains("file1"))
		assert.MainViewContent(Contains("+second line\n+third line"))
	},
})
