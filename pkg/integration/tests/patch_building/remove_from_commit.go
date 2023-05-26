package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoveFromCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a custom patch from a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content\n")
		shell.CreateFileAndAdd("file2", "file2 content\n")
		shell.Commit("first commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("first commit").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
				Contains("file2"),
			).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Views().PatchBuildingSecondary().Content(Contains("+file1 content"))

		t.Common().SelectPatchOption(Contains("Remove patch from original commit"))

		t.Views().Files().IsEmpty()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file2").IsSelected(),
			).
			PressEscape()

		t.Views().Main().
			Content(Contains("+file2 content"))

		t.Views().Commits().
			Lines(
				Contains("first commit").IsSelected(),
			)
	},
})
