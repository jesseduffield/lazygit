package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ApplyInReverseWithConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Apply a custom patch in reverse, resulting in a conflict",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content\n")
		shell.CreateFileAndAdd("file2", "file2 content\n")
		shell.Commit("first commit")
		shell.UpdateFileAndAdd("file1", "file1 content\nmore file1 content\n")
		shell.UpdateFileAndAdd("file2", "file2 content\nmore file2 content\n")
		shell.Commit("second commit")
		shell.UpdateFileAndAdd("file1", "file1 content\nmore file1 content\neven more file1\n")
		shell.Commit("third commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("third commit").IsSelected(),
				Contains("second commit"),
				Contains("first commit"),
			).
			NavigateToLine(Contains("second commit")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("M").Contains("file1").IsSelected(),
				Contains("M").Contains("file2"),
			).
			// Add both files to the patch; the first will conflict, the second won't
			PressPrimaryAction().
			Tap(func() {
				t.Views().Information().Content(Contains("Building patch"))

				t.Views().PatchBuildingSecondary().Content(
					Contains("+more file1 content"))
			}).
			SelectNextItem().
			PressPrimaryAction()

		t.Views().PatchBuildingSecondary().Content(
			Contains("+more file1 content").Contains("+more file2 content"))

		t.Common().SelectPatchOption(Contains("Apply patch in reverse"))

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Applied patch to 'file1' with conflicts.")).
			Confirm()

		t.Views().Files().
			Focus().
			Lines(
				Contains("UU").Contains("file1").IsSelected(),
			).
			PressPrimaryAction()

		t.Views().MergeConflicts().
			IsFocused().
			ContainsLines(
				Contains("file1 content"),
				Contains("<<<<<<< ours").IsSelected(),
				Contains("more file1 content").IsSelected(),
				Contains("even more file1").IsSelected(),
				Contains("=======").IsSelected(),
				Contains(">>>>>>> theirs"),
			).
			SelectNextItem().
			PressPrimaryAction()

		t.Views().Files().
			Focus().
			Lines(
				Contains("M").Contains("file1").IsSelected(),
				Contains("M").Contains("file2"),
			)

		t.Views().Main().
			ContainsLines(
				Contains(" file1 content"),
				Contains("-more file1 content"),
				Contains("-even more file1"),
			)
	},
})
