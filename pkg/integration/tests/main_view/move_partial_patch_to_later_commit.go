package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePartialPatchToLaterCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move one of two adjacent additions from a commit to a later commit",
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
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+1st line")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Commits().Focus().SelectPreviousItem()
		t.Common().SelectPatchOption(Contains("Move patch to selected commit"))

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
				Equals("▼ /").IsSelected(),
				Contains("file1"),
				Contains("unrelated-file"),
			).
			SelectNextItem().
			Tap(func() {
				t.Views().Main().Content(Contains("+1st line\n 2nd line"))
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
				t.Views().Main().Content(Contains("+2nd line").DoesNotContain("1st line"))
			})
	},
})
