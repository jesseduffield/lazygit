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

		// add the file to the patch
		input.PrimaryAction()

		input.PressKeys(keys.Universal.DiffingMenu)
		assert.InMenu()
		assert.MatchCurrentViewTitle(Equals("Diffing"))
		input.NavigateToListItemContainingText("exit diff mode")
		input.Confirm()

		assert.MatchViewContent("information", NotContains("building patch"))

		input.PressKeys(keys.Universal.CreatePatchOptionsMenu)
		assert.InMenu()
		assert.MatchCurrentViewTitle(Equals("Patch Options"))
		// including the keybinding 'a' here to distinguish the menu item from the 'apply patch in reverse' item
		input.NavigateToListItemContainingText("a apply patch")
		input.Confirm()

		input.SwitchToFilesWindow()

		assert.MatchSelectedLine(Contains("file1"))
		assert.MatchMainViewContent(Contains("+second line"))
	},
})
