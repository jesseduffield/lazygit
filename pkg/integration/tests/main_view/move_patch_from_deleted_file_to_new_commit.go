package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePatchFromDeletedFileToNewCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move part of a deleted file from a commit to a new commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "1st line\n2nd line\n3rd line\n")
		shell.Commit("first commit")
		shell.DeleteFileAndAdd("file1")
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
				Contains("D file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			NavigateToLine(Contains("-2nd line")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Common().SelectPatchOption(Contains("Move patch into new commit after the original commit"))

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("")).
			Type("new commit").Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("new commit").IsSelected(),
				Contains("commit to move from"),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("D file1").IsSelected(),
			).
			Tap(func() {
				t.Views().Main().ContainsLines(
					Equals("-2nd line"),
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
				Contains("M file1").IsSelected(),
			).
			Tap(func() {
				t.Views().Main().ContainsLines(
					Equals("-1st line"),
					Equals(" 2nd line"),
					Equals("-3rd line"),
				)
			})
	},
})
