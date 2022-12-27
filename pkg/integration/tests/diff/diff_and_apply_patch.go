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
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Views().Branches().
			Focus().
			Lines(
				Contains("branch-a"),
				Contains("branch-b"),
			).
			Press(keys.Universal.DiffingMenu)

		input.ExpectMenu().Title(Equals("Diffing")).Select(Equals("diff branch-a")).Confirm()

		input.Views().Information().Content(Contains("showing output for: git diff branch-a branch-a"))

		input.Views().Branches().
			IsFocused().
			SelectNextItem()

		input.Views().Information().Content(Contains("showing output for: git diff branch-a branch-b"))
		input.Views().Main().Content(Contains("+second line"))

		input.Views().Branches().
			PressEnter()

		input.Views().SubCommits().
			IsFocused().
			SelectedLine(Contains("update"))

		input.Views().Main().Content(Contains("+second line"))

		input.Views().SubCommits().
			PressEnter()

		input.Views().CommitFiles().
			IsFocused().
			SelectedLine(Contains("file1"))

		input.Views().Main().Content(Contains("+second line"))

		input.Views().CommitFiles().
			PressPrimaryAction(). // add the file to the patch
			Press(keys.Universal.DiffingMenu)

		input.ExpectMenu().Title(Equals("Diffing")).Select(Contains("exit diff mode")).Confirm()

		input.Views().Information().Content(DoesNotContain("building patch"))

		input.Views().CommitFiles().
			Press(keys.Universal.CreatePatchOptionsMenu)

		// adding the regex '$' here to distinguish the menu item from the 'apply patch in reverse' item
		input.ExpectMenu().Title(Equals("Patch Options")).Select(MatchesRegexp("apply patch$")).Confirm()

		input.Views().Files().
			Focus().
			SelectedLine(Contains("file1"))

		input.Views().Main().Content(Contains("+second line"))
	},
})
