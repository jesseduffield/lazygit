package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiffAndApplyPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a patch from the diff between two branches and apply the patch.",
	ExtraCmdArgs: []string{},
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
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch-a"),
				Contains("branch-b"),
			).
			Press(keys.Universal.DiffingMenu)

		t.ExpectPopup().Menu().Title(Equals("Diffing")).Select(Equals("Diff branch-a")).Confirm()

		t.Views().Information().Content(Contains("Showing output for: git diff branch-a branch-a"))

		t.Views().Branches().
			IsFocused().
			SelectNextItem().
			Tap(func() {
				t.Views().Information().Content(Contains("Showing output for: git diff branch-a branch-b"))
				t.Views().Main().Content(Contains("+second line"))
			}).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			SelectedLine(Contains("update")).
			Tap(func() {
				t.Views().Main().Content(Contains("+second line"))
			}).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			SelectedLine(Contains("file1")).
			Tap(func() {
				t.Views().Main().Content(Contains("+second line"))
			}).
			PressPrimaryAction(). // add the file to the patch
			Press(keys.Universal.DiffingMenu).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Diffing")).Select(Contains("Exit diff mode")).Confirm()

				t.Views().Information().Content(Contains("Building patch"))
			}).
			Press(keys.Universal.CreatePatchOptionsMenu)

		// adding the regex '$' here to distinguish the menu item from the 'Apply patch in reverse' item
		t.ExpectPopup().Menu().Title(Equals("Patch options")).Select(MatchesRegexp("Apply patch$")).Confirm()

		t.Views().Files().
			Focus().
			SelectedLine(Contains("file1"))

		t.Views().Main().Content(Contains("+second line"))
	},
})
