package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowAmendCommitMessageDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show a diff of the commit message when selecting an amend! commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommitWithBody("Fix the widget",
				"The frobnicator was not initialised properly at all.").
			EmptyCommitWithBody("amend! Fix the widget",
				"Fix the widget on startup\n\nThe frobnicator was not initialised at all.").
			EmptyCommitWithBody("amend! Fix the widget",
				"Fix the widget on startup\n\nThe frobnicator was not properly initialised at all.")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// git's default colors for added and removed lines
		green := "#008000"
		red := "#800000"

		t.Views().Commits().
			Focus().
			Lines(
				Contains("amend! Fix the widget").IsSelected(),
				Contains("amend! Fix the widget"),
				Contains("Fix the widget"),
			)

		// The message of the topmost amend! commit is compared with the message
		// of the amend! commit below it, not with the one of the commit they
		// both apply to.
		t.Views().Main().
			TopLines(
				Contains("Commit message changes compared to"),
				Equals(" Fix the widget on startup"),
				Equals(" "),
				Equals("-The frobnicator was not initialised at all."),
				Equals("+The frobnicator was not properly initialised at all."),
				Contains("───"),
			).
			ContainsColoredText(green, "+The frobnicator was not properly initialised at all.").
			ContainsColoredText(red, "-The frobnicator was not initialised at all.")

		t.Views().Commits().
			SelectNextItem()

		// The other amend! commit is compared with the commit it applies to.
		t.Views().Main().
			TopLines(
				Contains("Commit message changes compared to"),
				Equals("-Fix the widget"),
				Equals("+Fix the widget on startup"),
				Equals(" "),
				Equals("-The frobnicator was not initialised properly at all."),
				Equals("+The frobnicator was not initialised at all."),
				Contains("───"),
			)

		t.Views().Commits().
			SelectNextItem()

		// The commit that they both apply to gets no such header.
		t.Views().Main().
			TopLines(
				Contains("commit "),
			)
	},
})
