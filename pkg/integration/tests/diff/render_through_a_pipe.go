package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RenderThroughAPipe = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A stdin filter renders the diff that it is fed through a pipe",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// Reports the width it was given, then passes the diff through. git
		// only runs a filter of its own when it talks to a terminal, so the
		// filter running at all says the pipeline was built here.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Command: `echo "rendered at {{width}} columns"; cat`},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\n")
		shell.Commit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("one").IsSelected(),
			)

		t.Views().Main().
			// The width reaches the filter on its command line, since with no
			// terminal it has nowhere to read it from.
			Content(MatchesRegexp(`rendered at \d+ columns`)).
			ContainsLines(
				Equals("diff --git a/file1 b/file1"),
				Contains("new file mode"),
				Contains("index "),
				Equals("--- /dev/null"),
				Equals("+++ b/file1"),
				Equals("@@ -0,0 +1 @@"),
				Equals("+one"),
			)
	},
})
