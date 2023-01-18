package patch_building

import (
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BuildPatchAndCopyToClipboard = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a patch from the commits and copy the patch to clipbaord.",
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
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch-a").IsSelected(),
				Contains("branch-b"),
			).
			Press(keys.Universal.NextItem).
			PressEnter().
			PressEnter()
		t.Views().
			CommitFiles().
			Lines(
				Contains("M file1").IsSelected(),
			).
			PressPrimaryAction().Press(keys.Universal.CreatePatchOptionsMenu)

		t.ExpectPopup().Menu().Title(Equals("Patch Options")).Select(Contains("Copy patch to clipboard")).Confirm()

		text, err := clipboard.ReadAll()
		if err != nil {
			t.Fail(err.Error())
		}

		if !strings.HasPrefix(text, "diff --git a/file1 b/file1") {
			t.Fail("Text from clipboard did not match with git diff")
		}
	},
})
