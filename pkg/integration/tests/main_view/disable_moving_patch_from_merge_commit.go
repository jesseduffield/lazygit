package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisableMovingPatchFromMergeCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Moving a patch from a merge commit to another commit, or into a new commit before it, is disabled",
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
			Merge("feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("Merge branch 'feature' into main-branch").IsSelected(),
				Contains("feature commit"),
				Contains("main commit"),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file2").IsSelected(),
			).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("first commit"))

		t.GlobalPress(keys.Universal.CreatePatchOptionsMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Patch options")).
			Select(Contains("Move patch to selected commit")).
			Tooltip(Contains("Disabled: Not supported for a patch from a merge commit")).
			Select(Contains("Move patch into new commit before the original commit")).
			Tooltip(Contains("Disabled: Not supported for a patch from a merge commit")).
			Cancel()
	},
})
