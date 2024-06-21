package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemovePartsOfAddedFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a custom patch from a file that was added in a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")

		shell.CreateFileAndAdd("file1", "1st line\n2nd line\n3rd line\n")
		shell.Commit("commit to remove from")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit to remove from").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("A file1").IsSelected(),
			).
			PressEnter()

		t.Views().PatchBuilding().
			IsFocused().
			SelectNextItem().
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Common().SelectPatchOption(Contains("Remove patch from original commit"))

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("A file1").IsSelected(),
			).
			PressEscape()

		t.Views().Main().ContainsLines(
			Equals("+1st line"),
			Equals("+3rd line"),
		)
	},
})
