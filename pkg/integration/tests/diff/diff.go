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
		input.SwitchToBranchesView()

		assert.Views().Current().TopLines(
			Contains("branch-a"),
			Contains("branch-b"),
		)
		input.Press(keys.Universal.DiffingMenu)
		input.Menu().Title(Equals("Diffing")).Select(Contains(`diff branch-a`)).Confirm()

		assert.Views().Current().Name("localBranches")

		assert.Views().ByName("information").Content(Contains("showing output for: git diff branch-a branch-a"))
		input.NextItem()
		assert.Views().ByName("information").Content(Contains("showing output for: git diff branch-a branch-b"))
		assert.Views().Main().Content(Contains("+second line"))

		input.Enter()
		assert.Views().Current().Name("subCommits")
		assert.Views().Main().Content(Contains("+second line"))
		assert.Views().Current().SelectedLine(Contains("update"))
		input.Enter()
		assert.Views().Current().Name("commitFiles").SelectedLine(Contains("file1"))
		assert.Views().Main().Content(Contains("+second line"))

		input.Press(keys.Universal.Return)
		input.Press(keys.Universal.Return)
		assert.Views().Current().Name("localBranches")

		input.Press(keys.Universal.DiffingMenu)
		input.Menu().Title(Equals("Diffing")).Select(Contains("reverse diff direction")).Confirm()
		assert.Views().ByName("information").Content(Contains("showing output for: git diff branch-a branch-b -R"))
		assert.Views().Main().Content(Contains("-second line"))
	},
})
