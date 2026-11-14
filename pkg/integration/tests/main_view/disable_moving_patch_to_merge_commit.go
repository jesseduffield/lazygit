package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisableMovingPatchToMergeCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Moving a patch to an earlier or later merge commit is disabled",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("main-branch").
			CreateFileAndAdd("file1", "file1 content\n").
			Commit("first commit").
			NewBranch("feature").
			CreateFileAndAdd("file2", "file2 content\n").
			Commit("feature commit").
			Checkout("main-branch").
			CreateFileAndAdd("file3", "file3 content\n").
			Commit("main commit").
			Merge("feature").
			CreateFileAndAdd("file4", "file4 content\n").
			Commit("top commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		checkMovingPatchToMergeCommitIsDisabled := func(commitToMoveFrom string, file string) {
			t.Views().Commits().
				Focus().
				NavigateToLine(Contains(commitToMoveFrom)).
				PressEnter()

			t.Views().CommitFiles().
				IsFocused().
				Lines(
					Contains(file).IsSelected(),
				).
				PressPrimaryAction()

			t.Views().Information().Content(Contains("Building patch"))

			t.Views().Commits().
				Focus().
				NavigateToLine(Contains("Merge branch 'feature' into main-branch"))

			t.GlobalPress(keys.Universal.CreatePatchOptionsMenu)

			t.ExpectPopup().Menu().
				Title(Equals("Patch options")).
				Select(Contains("Move patch to selected commit")).
				Tooltip(Contains("Disabled: Cannot move a patch to a merge commit")).
				Cancel()

			t.Common().ResetCustomPatch()
		}

		t.Views().Commits().
			Lines(
				Contains("top commit").IsSelected(),
				Contains("Merge branch 'feature' into main-branch"),
				Contains("feature commit"),
				Contains("main commit"),
				Contains("first commit"),
			)

		checkMovingPatchToMergeCommitIsDisabled("top commit", "file4")
		checkMovingPatchToMergeCommitIsDisabled("main commit", "file3")
	},
})
