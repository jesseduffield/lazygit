package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardAllChangesOfSquashMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard all changes of a squash merge that was left uncommitted, which ends the squash merge",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("original-branch").
			CreateFileAndAdd("file", "original\n").
			Commit("one").
			NewBranch("feature-branch").
			UpdateFileAndAdd("file", "feature\n").
			CreateFileAndAdd("feature-file", "content").
			Commit("feature").
			Checkout("original-branch").
			CreateFile("untracked-file", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("original-branch").IsSelected(),
				Contains("feature-branch"),
			).
			SelectNextItem().
			Press(keys.Branches.MergeIntoCurrentBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Squash merge and leave uncommitted")).
			Confirm()

		t.Views().Status().Content(Equals("(merging) repo → original-branch"))

		t.Views().Files().
			Focus().
			Lines(
				Equals("▼ /"),
				Equals("  A  feature-file"),
				Equals("  M  file"),
				Equals("  ?? untracked-file"),
			).
			NavigateToLine(Equals("  M  file")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			Lines(
				Equals("▼ /"),
				Equals("  A  feature-file"),
				Equals("  ?? untracked-file"),
			).
			NavigateToLine(Equals("  A  feature-file"))

		// There are still changes left, so the squash merge is still in progress
		t.Views().Status().Content(Equals("(merging) repo → original-branch"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			Lines(
				Equals("?? untracked-file").IsSelected(),
			)

		t.Views().Status().Content(Equals("repo → original-branch"))

		t.FileSystem().PathNotPresent(".git/SQUASH_MSG")
	},
})
