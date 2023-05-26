package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CopyPatchToClipboard = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a patch from the commits and copy the patch to clipbaord.",
	ExtraCmdArgs: []string{},
	Skip:         true, // skipping because CI doesn't have clipboard functionality
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
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Common().SelectPatchOption(Contains("copy patch to clipboard"))

		t.ExpectToast(Contains("Patch copied to clipboard"))

		t.ExpectClipboard(Contains("diff --git a/file1 b/file1"))
	},
})
