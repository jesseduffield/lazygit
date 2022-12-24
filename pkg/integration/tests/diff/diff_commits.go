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

		input.PressKeys(keys.Universal.DiffingMenu)
		assert.InMenu()
		assert.CurrentViewTitle(Equals("Diffing"))
		assert.SelectedLine(Contains("diff"))
		input.Confirm()
		assert.NotInPopup()

		assert.ViewContent("information", Contains("showing output for: git diff"))

		input.NextItem()
		input.NextItem()

		assert.SelectedLine(Contains("first commit"))

		assert.MainViewContent(Contains("-second line\n-third line"))

		input.PressKeys(keys.Universal.DiffingMenu)
		assert.InMenu()
		input.NavigateToListItemContainingText("reverse diff direction")
		input.Confirm()

		assert.MainViewContent(Contains("+second line\n+third line"))

		input.Enter()

		assert.CurrentViewName("commitFiles")
		assert.SelectedLine(Contains("file1"))
		assert.MainViewContent(Contains("+second line\n+third line"))
	},
})
