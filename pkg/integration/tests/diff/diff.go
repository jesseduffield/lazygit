package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Diff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "View the diff of two branches, then view the reverse diff",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "first line")
		shell.Commit("first commit")

		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "first line\nsecond line")
		shell.Commit("update")

		shell.Checkout("branch-a")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToBranchesWindow()
		assert.CurrentViewName("localBranches")

		assert.MatchSelectedLine(Contains("branch-a"))
		input.PressKeys(keys.Universal.DiffingMenu)
		assert.InMenu()
		assert.MatchCurrentViewTitle(Equals("Diffing"))
		assert.MatchSelectedLine(Contains("diff branch-a"))
		input.Confirm()

		assert.CurrentViewName("localBranches")

		assert.MatchViewContent("information", Contains("showing output for: git diff branch-a branch-a"))
		input.NextItem()
		assert.MatchViewContent("information", Contains("showing output for: git diff branch-a branch-b"))
		assert.MatchMainViewContent(Contains("+second line"))

		input.Enter()
		assert.CurrentViewName("subCommits")
		assert.MatchMainViewContent(Contains("+second line"))
		assert.MatchSelectedLine(Contains("update"))
		input.Enter()
		assert.CurrentViewName("commitFiles")
		assert.MatchSelectedLine(Contains("file1"))
		assert.MatchMainViewContent(Contains("+second line"))

		input.PressKeys(keys.Universal.Return)
		input.PressKeys(keys.Universal.Return)
		assert.CurrentViewName("localBranches")

		input.PressKeys(keys.Universal.DiffingMenu)
		assert.InMenu()
		input.NavigateToListItemContainingText("reverse diff direction")
		input.Confirm()
		assert.MatchViewContent("information", Contains("showing output for: git diff branch-a branch-b -R"))
		assert.MatchMainViewContent(Contains("-second line"))
	},
})
