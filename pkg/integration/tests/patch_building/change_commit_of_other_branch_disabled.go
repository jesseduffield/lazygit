package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChangeCommitOfOtherBranchDisabled = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With the whole graph shown, neither discarding file changes nor custom patches can change a commit of another branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.ShowWholeGraph = true
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("base.txt", "base\n").
			Commit("base").
			NewBranch("dev-a").
			CreateFileAndAdd("a.txt", "a\n").
			Commit("dev-a-1").
			NewBranchFrom("dev-b", "master").
			CreateFileAndAdd("b.txt", "b\n").
			Commit("dev-b-1")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		reason := "Only commits of the current branch can be changed; the others are only listed because the whole git graph is shown"

		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("dev-a-1")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("A a.txt").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Equals(reason)).
					Confirm()
			}).
			PressPrimaryAction()

		t.GlobalPress(keys.Universal.CreatePatchOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Patch options")).
			Lines(
				Contains("Reset patch"),
				Contains("Apply patch"),
				Contains("Apply patch in reverse"),
				Contains("Copy patch to clipboard"),
				Contains("Cancel"),
			).
			Select(Contains("Reset patch")).
			Confirm()

		t.Views().CommitFiles().
			PressEscape()

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("dev-b-1")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("A b.txt").IsSelected(),
			).
			PressPrimaryAction().
			PressEscape()

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("dev-a-1"))

		t.GlobalPress(keys.Universal.CreatePatchOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Patch options")).
			Select(Contains("Move patch to selected commit")).
			Confirm()

		t.ExpectToast(Equals("Disabled: " + reason))
	},
})
