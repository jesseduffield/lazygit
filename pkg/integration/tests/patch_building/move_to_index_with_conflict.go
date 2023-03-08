package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToIndexWithConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit to the index, causing a conflict",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "file1 content with old changes")
		shell.Commit("second commit")

		shell.UpdateFileAndAdd("file1", "file1 content with new changes")
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
			SelectNextItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("building patch"))

		t.Common().SelectPatchOption(Contains("move patch out into index"))

		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU").Contains("file1"),
			).
			PressPrimaryAction()

		t.Views().MergeConflicts().
			IsFocused().
			ContainsLines(
				Contains("<<<<<<< HEAD").IsSelected(),
				Contains("file1 content").IsSelected(),
				Contains("=======").IsSelected(),
				Contains("file1 content with new changes"),
				Contains(">>>>>>>"),
			).
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
				Contains("<<<<<<< ours"),
				Contains("file1 content"),
				Contains("======="),
				Contains("file1 content with old changes"),
				Contains(">>>>>>> theirs"),
			).
			IsFocused()
	},
})
