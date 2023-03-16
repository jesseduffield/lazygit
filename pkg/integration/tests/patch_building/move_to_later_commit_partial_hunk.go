package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToLaterCommitPartialHunk = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit to a later commit, with only parts of a hunk in the patch",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "1st line\n2nd line\n")
		shell.Commit("commit to move from")

		shell.UpdateFileAndAdd("unrelated-file", "")
		shell.Commit("destination commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("destination commit").IsSelected(),
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
			PressEnter().
			PressPrimaryAction().
			PressEscape()

		t.Views().Information().Content(Contains("building patch"))

		t.Views().CommitFiles().
			IsFocused().
			PressEscape()

		t.Views().Commits().
			IsFocused().
			SelectPreviousItem()

		t.Common().SelectPatchOption(Contains("move patch to selected commit"))

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("destination commit").IsSelected(),
				Contains("commit to move from"),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
				Contains("unrelated-file"),
			).
			Tap(func() {
				t.Views().Main().
					Content(Contains("+1st line\n 2nd line"))
			}).
			PressEscape()

		t.Views().Commits().
			IsFocused().
			SelectNextItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Tap(func() {
				t.Views().Main().
					Content(Contains("+2nd line").
						DoesNotContain("1st line"))
			})
	},
})
