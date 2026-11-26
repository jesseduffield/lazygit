package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AbortCherryPickAfterEditingPickedFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Abort a cherry-pick that stopped with conflicts after editing a file that it picked without conflicts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("conflicting-file", "original\n").
			CreateFileAndAdd("picked-file", "original\n").
			Commit("original").
			NewBranch("feature").
			UpdateFileAndAdd("conflicting-file", "feature\n").
			UpdateFileAndAdd("picked-file", "feature\n").
			Commit("feature change").
			Checkout("master").
			UpdateFileAndAdd("conflicting-file", "master\n").
			Commit("master change").
			RunCommandExpectError([]string{"git", "cherry-pick", "feature"}).
			UpdateFile("picked-file", "feature\nfix the build\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("(cherry-picking) repo → master"))

		t.Views().Files().
			IsFocused().
			Lines(
				Equals("UU conflicting-file").IsSelected(),
			)

		t.GlobalPress(keys.Universal.CreateRebaseOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Cherry-pick options")).
			Select(Contains("abort")).
			Confirm()

		t.ExpectPopup().Confirmation().
			Title(Equals("Abort the cherry-pick")).
			Content(Equals("Aborting the cherry-pick also discards the unstaged changes of these files. Are you sure?\n\n- picked-file")).
			Confirm()

		t.Views().Status().Content(Equals("repo → master"))

		t.Views().Files().IsEmpty()

		t.Views().Commits().
			Lines(
				Contains("master change"),
				Contains("original"),
			)
	},
})
