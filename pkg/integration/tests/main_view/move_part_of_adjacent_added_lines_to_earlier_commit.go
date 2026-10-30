package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePartOfAdjacentAddedLinesToEarlierCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move only one of two adjacent added lines from a commit to an earlier commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "")
		shell.Commit("first commit")

		shell.CreateFileAndAdd("unrelated-file", "")
		shell.Commit("destination commit")

		shell.UpdateFileAndAdd("file1", "1st line\n2nd line\n")
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
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+1st line")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Commits().Focus().SelectNextItem()
		t.Common().SelectPatchOption(Contains("Move patch to selected commit"))

		/* EXPECTED:
		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("commit to move from"),
				Contains("destination commit").IsSelected(),
				Contains("first commit"),
			).
			Tap(func() {
				t.Views().Main().Content(Contains("+1st line").DoesNotContain("2nd line"))
			}).
			SelectPreviousItem().
			Tap(func() {
				t.Views().Main().Content(Contains(" 1st line\n+2nd line"))
			})
		ACTUAL: */
		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU").Contains("file1"),
			)
	},
})
