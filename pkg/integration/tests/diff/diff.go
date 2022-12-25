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

		assert.CurrentViewTopLines(
			Contains("branch-a"),
			Contains("branch-b"),
		)
		input.Press(keys.Universal.DiffingMenu)
		input.Menu(Equals("Diffing"), Contains(`diff branch-a`))

		assert.CurrentViewName("localBranches")

		assert.ViewContent("information", Contains("showing output for: git diff branch-a branch-a"))
		input.NextItem()
		assert.ViewContent("information", Contains("showing output for: git diff branch-a branch-b"))
		assert.MainViewContent(Contains("+second line"))

		input.Enter()
		assert.CurrentViewName("subCommits")
		assert.MainViewContent(Contains("+second line"))
		assert.SelectedLine(Contains("update"))
		input.Enter()
		assert.CurrentViewName("commitFiles")
		assert.SelectedLine(Contains("file1"))
		assert.MainViewContent(Contains("+second line"))

		input.Press(keys.Universal.Return)
		input.Press(keys.Universal.Return)
		assert.CurrentViewName("localBranches")

		input.Press(keys.Universal.DiffingMenu)
		input.Menu(Equals("Diffing"), Contains("reverse diff direction"))
		assert.ViewContent("information", Contains("showing output for: git diff branch-a branch-b -R"))
		assert.MainViewContent(Contains("-second line"))
	},
})
