package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToNewCommitPartialHunk = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit to a new commit, with only parts of a hunk in the patch",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "1st line\n2nd line\n")
		shell.Commit("commit to move from")

		shell.UpdateFileAndAdd("file1", "1st line\n2nd line\n3rd line\n")
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
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().PatchBuilding().
			IsFocused().
			PressEnter().
			PressPrimaryAction()

		t.Views().Information().Content(Contains("building patch"))

		t.Common().SelectPatchOption(Contains("move patch into new commit"))

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Tap(func() {
				t.Views().Main().
					Content(Contains("+1st line\n 2nd line"))
			}).
			PressEscape()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("third commit"),
				Contains(`Split from "commit to move from"`).IsSelected(),
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
			Tap(func() {
				t.Views().Main().
					Content(Contains("+2nd line").
						DoesNotContain("1st line"))
			})
	},
})
