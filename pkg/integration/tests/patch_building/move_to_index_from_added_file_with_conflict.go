package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToIndexFromAddedFileWithConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a file that was added in a commit to the index, causing a conflict",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")

		shell.CreateFileAndAdd("file1", "1st line\n2nd line\n3rd line\n")
		shell.Commit("commit to move from")
		shell.UpdateFileAndAdd("file1", "1st line\n2nd line changed\n3rd line\n")
		shell.Commit("conflicting change")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("conflicting change").IsSelected(),
				Contains("commit to move from"),
				Contains("first commit"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().PatchBuilding().
			IsFocused().
			SelectNextItem().
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Common().SelectPatchOption(Contains("Move patch out into index"))

		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU").Contains("file1"),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			ContainsLines(
				Contains("1st line"),
				Contains("<<<<<<< HEAD"),
				Contains("======="),
				Contains("2nd line changed"),
				Contains(">>>>>>>"),
				Contains("3rd line"),
			).
			SelectNextItem().
			PressPrimaryAction()

		t.Common().ContinueOnConflictsResolved()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Applied patch to 'file1' with conflicts")).
			Confirm()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU").Contains("file1"),
			).
			PressEnter()

		t.Views().MergeConflicts().
			TopLines(
				Contains("1st line"),
				Contains("<<<<<<< ours"),
				Contains("2nd line changed"),
				Contains("======="),
				Contains("2nd line"),
				Contains(">>>>>>> theirs"),
				Contains("3rd line"),
			).
			IsFocused()
	},
})
