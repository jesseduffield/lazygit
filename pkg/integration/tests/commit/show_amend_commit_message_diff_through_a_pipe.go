package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowAmendCommitMessageDiffThroughAPipe = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show the commit message diff of an amend! commit through a diff renderer that is fed through a pipe rather than run in a pty",
	ExtraCmdArgs: []string{},
	Skip:         false,
	// This is how a render works on Windows, where a pty can't carry a
	// renderer's output faithfully. Ask for it here so that the path is
	// covered on the platforms the integration tests do run on.
	ExtraEnvVars: map[string]string{"LAZYGIT_RENDER_WITHOUT_PTY": "1"},
	SetupConfig: func(cfg *config.AppConfig) {
		// Says so if it is writing into a pipe, then passes the diff through.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Command: `[ -t 1 ] || echo "rendered into a pipe"; cat`},
		}
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

		t.Views().Main().
			TopLines(
				Contains("Commit message changes compared to"),
				Equals("rendered into a pipe"),
				Equals("diff --git old message new message"),
			).
			ContainsLines(
				Equals("-Fix the widget"),
				Equals("+Fix the widget on startup"),
			)
	},
})
