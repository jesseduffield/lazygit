package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiffRendererMetadata = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A diff renderer is told that we understand the OSC 1717 metadata protocol, and the records it emits don't show up in the rendered diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// A fake conforming renderer: it announces the protocol with a
		// version-only record, reports the protocol versions it was offered, and
		// then passes the diff through with a per-line record in front of every
		// line.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Command: `printf '\033]1717;1\007'; ` +
				`printf 'OFFERED:%s\n' "$OSC1717"; ` +
				`while IFS= read -r line; do printf '\033]1717;1;c;1;;file1\007%s\n' "$line"; done`},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("one").IsSelected(),
			)

		t.Views().Main().
			// The renderer was offered the protocol version we understand.
			Content(Contains("OFFERED:V1")).
			// Its records are escape sequences, so none of them reaches the
			// screen; the diff reads exactly as the renderer wrote it.
			ContainsLines(
				Equals("diff --git a/file1 b/file1"),
				Contains("new file mode"),
				Contains("index "),
				Equals("--- /dev/null"),
				Equals("+++ b/file1"),
				Equals("@@ -0,0 +1,2 @@"),
				Equals("+one"),
				Equals("+two"),
			)
	},
})
