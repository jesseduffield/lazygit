package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToNewCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit to a new commit",
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

		t.Common().SelectPatchOption(Contains("move patch into new commit"))

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEscape()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("third commit"),
				Contains(`Split from "second commit"`).IsSelected(),
				Contains("second commit"),
				Contains("first commit"),
			).
			SelectNextItem().
			PressEnter()

		// the original commit has no more files in it
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("(none)"),
			)
	},
})
