package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveAmbiguouslyPlacedAddedLinesToNewCommitBefore = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move some added lines from a commit to a new commit before it, where a diff could also place those lines one line further down",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "first line\n}\nlast line\n")
		shell.Commit("first commit")

		// A diff can show the added "}" and "added line" as "added line" and "}"
		// added one line further down, after the existing "}"
		shell.UpdateFileAndAdd("file1", "first line\nextra line\n}\nadded line\n}\nlast line\n")
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
			SelectedLines(Contains("+extra line")).
			NavigateToLine(Contains("+}")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+added line")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Common().SelectPatchOption(Contains("Move patch into new commit before the original commit"))

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("")).
			Type("new commit").Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("commit to move from").IsSelected(),
				Contains("new commit"),
				Contains("first commit"),
			)

		t.FileSystem().FileContent(
			"file1",
			/* EXPECTED:
			Equals("first line\nextra line\n}\nadded line\n}\nlast line\n"),
			ACTUAL: */
			Equals("first line\nextra line\n}\nadded line\n}\nadded line\n}\nlast line\n"),
		)
	},
})
