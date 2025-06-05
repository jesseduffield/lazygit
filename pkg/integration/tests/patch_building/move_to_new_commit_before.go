package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToNewCommitBefore = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit to a new commit before the original one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.26.0"),
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateFileAndAdd("dir/file1", "file1 content")
		shell.CreateFileAndAdd("dir/file2", "file2 content")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("dir/file1", "file1 content with old changes")
		shell.DeleteFileAndAdd("dir/file2")
		shell.CreateFileAndAdd("dir/file3", "file3 content")
		shell.Commit("commit to move from")

		shell.UpdateFileAndAdd("dir/file1", "file1 content with new changes")
		shell.Commit("third commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("third commit").IsSelected(),
				Contains("commit to move from"),
				Contains("first commit"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir").IsSelected(),
				Contains("  M file1"),
				Contains("  D file2"),
				Contains("  A file3"),
			).
			PressPrimaryAction().
			PressEscape()

		t.Views().Information().Content(Contains("Building patch"))

		t.Common().SelectPatchOption(Contains("Move patch into new commit before the original commit"))

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("")).
			Type("new commit").Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("third commit"),
				Contains("commit to move from").IsSelected(),
				Contains("new commit"),
				Contains("first commit"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir").IsSelected(),
				Contains("  M file1"),
				Contains("  D file2"),
				Contains("  A file3"),
			).
			PressEscape()

		t.Views().Commits().
			IsFocused().
			SelectPreviousItem().
			PressEnter()

		// the original commit has no more files in it
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("(none)"),
			)
	},
})
