package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToEarlierCommitFromAddedFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a file that was added in a commit to an earlier commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.EmptyCommit("destination commit")
		shell.CreateFileAndAdd("file1", "1st line\n2nd line\n3rd line\n")
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
				Contains("A file").IsSelected(),
			).
			PressEnter()

		t.Views().PatchBuilding().
			IsFocused().
			SelectNextItem().
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Views().Commits().
			Focus().
			SelectNextItem()

		t.Common().SelectPatchOption(Contains("Move patch to selected commit"))

		// This results in a conflict at the commit we're moving from, because
		// it tries to add a file that already exists
		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("AA").Contains("file"),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			TopLines(
				Contains("<<<<<<< HEAD"),
				Contains("2nd line"),
				Contains("======="),
				Contains("1st line"),
				Contains("2nd line"),
				Contains("3rd line"),
				Contains(">>>>>>>"),
			).
			SelectNextItem().
			PressPrimaryAction() // choose the version with all three lines

		t.Common().ContinueOnConflictsResolved()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit to move from"),
				Contains("destination commit").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("A file").IsSelected(),
			).
			Tap(func() {
				t.Views().Main().ContainsLines(
					Equals("+2nd line"),
				)
			}).
			PressEscape()

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("commit to move from")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("M file").IsSelected(),
			).
			Tap(func() {
				t.Views().Main().ContainsLines(
					Equals("+1st line"),
					Equals(" 2nd line"),
					Equals("+3rd line"),
				)
			})
	},
})
