package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AbortMergeAfterEditingMergedFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Abort a merge that stopped with conflicts after editing a file that it merged without conflicts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("conflicting-file", "original\n").
			CreateFileAndAdd("merged-file", "original\n").
			CreateFileAndAdd("unrelated-file", "original\n").
			Commit("original").
			NewBranch("feature").
			UpdateFileAndAdd("conflicting-file", "feature\n").
			UpdateFileAndAdd("merged-file", "feature\n").
			Commit("feature change").
			Checkout("master").
			UpdateFileAndAdd("conflicting-file", "master\n").
			Commit("master change").
			UpdateFile("unrelated-file", "unrelated change\n").
			RunCommandExpectError([]string{"git", "merge", "--no-edit", "feature"}).
			UpdateFile("merged-file", "feature\nfix the build\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("(merging) repo → master"))

		t.Views().Files().
			IsFocused().
			Lines(
				Equals("UU conflicting-file").IsSelected(),
			)

		t.GlobalPress(keys.Universal.CreateRebaseOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Merge options")).
			Select(Contains("abort")).
			Confirm()

		t.ExpectPopup().Confirmation().
			Title(Equals("Abort the merge")).
			Content(Equals("Aborting the merge also discards the unstaged changes of these files. Are you sure?\n\n- merged-file")).
			Confirm()

		t.Views().Status().Content(Equals("repo → master"))

		t.Views().Files().
			Lines(
				Equals(" M unrelated-file"),
			)

		t.FileSystem().FileContent("merged-file", Equals("original\n"))
	},
})
