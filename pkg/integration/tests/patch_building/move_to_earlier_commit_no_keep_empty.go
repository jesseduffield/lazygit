package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToEarlierCommitNoKeepEmpty = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit to an earlier commit, for older git versions that don't keep the empty commit",
	ExtraCmdArgs: "",
	Skip:         false,
	GitVersion:   Before("2.26.0"),
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateFileAndAdd("dir/file1", "file1 content")
		shell.CreateFileAndAdd("dir/file2", "file2 content")
		shell.Commit("first commit")

		shell.CreateFileAndAdd("unrelated-file", "")
		shell.Commit("destination commit")

		shell.UpdateFileAndAdd("dir/file1", "file1 content with old changes")
		shell.DeleteFileAndAdd("dir/file2")
		shell.CreateFileAndAdd("dir/file3", "file3 content")
		shell.Commit("commit to move from")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit to move from").IsSelected(),
				Contains("destination commit"),
				Contains("first commit"),
			).
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

		t.Views().Information().Content(Contains("building patch"))

		t.Views().Commits().
			IsFocused().
			SelectNextItem()

		t.Common().SelectPatchOption(Contains("move patch to selected commit"))

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("destination commit"),
				Contains("first commit").IsSelected(),
			).
			SelectPreviousItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir").IsSelected(),
				Contains("  M file1"),
				Contains("  D file2"),
				Contains("  A file3"),
				Contains("A unrelated-file"),
			).
			PressEscape()
	},
})
