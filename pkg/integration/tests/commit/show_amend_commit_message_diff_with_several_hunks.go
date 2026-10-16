package commit

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

const showAmendCommitMessageDiffWithSeveralHunksBody = `The frobnicator was not initialised.

It has to be initialised before the widget
draws itself for the first time, or the
widget shows garbage until it is resized.

This only showed on startup, so nobody
noticed it for a long time.

Fixes #123.`

var ShowAmendCommitMessageDiffWithSeveralHunks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The hunks of the commit message diff of an amend! commit are not headed by a line of the message",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommitWithBody("Fix the widget", showAmendCommitMessageDiffWithSeveralHunksBody).
			EmptyCommitWithBody("amend! Fix the widget",
				"Fix the widget on startup\n\n"+
					strings.Replace(showAmendCommitMessageDiffWithSeveralHunksBody, "#123", "#124", 1))
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("amend! Fix the widget").IsSelected(),
				Contains("Fix the widget"),
			)

		// git would head the second hunk with the last line above it that
		// starts with a letter, as if it named the function the hunk is in.
		t.Views().Main().
			TopLines(
				Contains("Commit message changes compared to"),
				Equals("-Fix the widget"),
				Equals("+Fix the widget on startup"),
				Equals(" "),
				Equals(" The frobnicator was not initialised."),
				Equals(" "),
				Equals("@@ -9,4 +9,4 @@"),
				Equals(" This only showed on startup, so nobody"),
				Equals(" noticed it for a long time."),
				Equals(" "),
				Equals("-Fixes #123."),
				Equals("+Fixes #124."),
				Contains("───"),
			)
	},
})
