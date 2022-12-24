package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiffAndApplyPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a patch from the diff between two branches and apply the patch.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "first line\n")
		shell.Commit("first commit")

		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "first line\nsecond line\n")
		shell.Commit("update")

		shell.Checkout("branch-a")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToBranchesWindow()
		assert.CurrentViewName("localBranches")

		assert.SelectedLine(Contains("branch-a"))
		input.Press(keys.Universal.DiffingMenu)

		input.Menu(Equals("Diffing"), Equals("diff branch-a"))

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

		// add the file to the patch
		input.PrimaryAction()

		input.Press(keys.Universal.DiffingMenu)
		input.Menu(Equals("Diffing"), Contains("exit diff mode"))

		assert.ViewContent("information", NotContains("building patch"))

		input.Press(keys.Universal.CreatePatchOptionsMenu)
		// adding the regex '$' here to distinguish the menu item from the 'apply patch in reverse' item
		input.Menu(Equals("Patch Options"), MatchesRegexp("apply patch$"))

		input.SwitchToFilesWindow()

		assert.SelectedLine(Contains("file1"))
		assert.MainViewContent(Contains("+second line"))
	},
})
