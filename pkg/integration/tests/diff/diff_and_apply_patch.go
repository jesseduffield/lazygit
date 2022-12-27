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
		input.SwitchToBranchesView()
		assert.Views().Current().Lines(
			Contains("branch-a"),
			Contains("branch-b"),
		)

		input.Press(keys.Universal.DiffingMenu)

		input.Menu().Title(Equals("Diffing")).Select(Equals("diff branch-a")).Confirm()

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
		assert.Views().Current().Name("commitFiles")
		assert.Views().Current().SelectedLine(Contains("file1"))
		assert.Views().Main().Content(Contains("+second line"))

		// add the file to the patch
		input.PrimaryAction()

		input.Press(keys.Universal.DiffingMenu)
		input.Menu().Title(Equals("Diffing")).Select(Contains("exit diff mode")).Confirm()

		assert.Views().ByName("information").Content(DoesNotContain("building patch"))

		input.Press(keys.Universal.CreatePatchOptionsMenu)
		// adding the regex '$' here to distinguish the menu item from the 'apply patch in reverse' item
		input.Menu().Title(Equals("Patch Options")).Select(MatchesRegexp("apply patch$")).Confirm()

		input.SwitchToFilesView()

		assert.Views().Current().SelectedLine(Contains("file1"))
		assert.Views().Main().Content(Contains("+second line"))
	},
})
