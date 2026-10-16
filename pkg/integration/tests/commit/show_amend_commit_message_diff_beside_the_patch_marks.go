package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowAmendCommitMessageDiffBesideThePatchMarks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The commit message diff of an amend! commit is laid out to the width that the custom patch's marks leave the diff below it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "one\n")
		shell.Commit("Fix the widget")
		shell.UpdateFileAndAdd("file", "one\ntwo\n")
		shell.Commit("amend! Fix the widget\n\nFix the widget on startup")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("amend! Fix the widget").IsSelected(),
				Contains("Fix the widget"),
			).
			Press(keys.Universal.FocusMainView)

		// The rule below the message diff is as wide as the diff is laid out to,
		// so it only fits on one line beside the marks if the message diff was
		// laid out to the width they leave.
		t.Views().Main().
			IsFocused().
			PressPrimaryAction().
			MarkedLines(
				Contains("+two"),
			).
			ContainsViewLines(
				Equals("+Fix the widget on startup"),
				MatchesRegexp("^─+$"),
				Contains("commit "),
			)
	},
})
