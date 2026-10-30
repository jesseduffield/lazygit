package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePatchFromAddedFileToEarlierCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move part of an added file from a commit to an earlier commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
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
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			NavigateToLine(Contains("+2nd line")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Commits().Focus().SelectNextItem()
		t.Common().SelectPatchOption(Contains("Move patch to selected commit"))

		t.Views().Commits().
			IsFocused().
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
