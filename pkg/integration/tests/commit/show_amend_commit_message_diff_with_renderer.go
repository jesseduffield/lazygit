package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowAmendCommitMessageDiffWithRenderer = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show the commit message diff of an amend! commit through a diff renderer",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// cat does nothing to the diff, but it is a diff renderer as far as
		// lazygit is concerned, so the diff takes the same route through it as
		// it would for a real one.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{{Command: "cat"}}
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommitWithBody("Fix the widget",
				"The frobnicator was not initialised.").
			EmptyCommitWithBody("amend! Fix the widget",
				"Fix the widget on startup\n\nThe frobnicator was not initialised.")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("amend! Fix the widget").IsSelected(),
				Contains("Fix the widget"),
			)

		// A diff renderer states the two sides of the diff in its own way, so
		// unlike for git's own diff, the header naming them is kept.
		t.Views().Main().
			TopLines(
				Contains("Commit message changes compared to"),
				Equals("diff --git old message new message"),
			).
			ContainsLines(
				Contains("--- old message"),
				Contains("+++ new message"),
				Contains("@@"),
				Equals("-Fix the widget"),
				Equals("+Fix the widget on startup"),
			)
	},
})
