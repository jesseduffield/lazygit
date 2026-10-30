package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePartOfAdjacentAddedLinesToNewCommitBefore = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move only one of two adjacent added lines from a commit to a new commit before it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "1st line\n2nd line\n")
		shell.Commit("commit to move from")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit to move from").IsSelected(),
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
		t.Common().SelectPatchOption(Contains("Move patch into new commit before the original commit"))

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("")).
			Type("new commit").Confirm()

		t.Views().Commits().
			IsFocused().
			/* EXPECTED:
			Lines(
				Contains("commit to move from").IsSelected(),
				Contains("new commit"),
				Contains("first commit"),
			).
			Tap(func() {
				t.Views().Main().Content(Contains(" 1st line\n+2nd line"))
			}).
			SelectNextItem().
			Tap(func() {
				t.Views().Main().Content(Contains("+1st line").DoesNotContain("2nd line"))
			})
			ACTUAL: */
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("pick").Contains("<-- CONFLICT --- commit to move from"),
				Contains("─── Commits"),
				Contains("new commit"),
				Contains("first commit"),
			)
	},
})
