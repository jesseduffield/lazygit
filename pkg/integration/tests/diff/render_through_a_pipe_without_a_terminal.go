package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RenderThroughAPipeWithoutATerminal = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A diff renderer fed through a pipe has no terminal to ask for its size",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// Reports the size of the terminal it can reach, the way a renderer that
		// lays itself out to the terminal's width (difftastic does) finds its width,
		// then passes the diff through.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Command: `echo "terminal: $(stty size </dev/tty 2>/dev/null || echo none)"; cat`},
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

		// The terminal it would find is lazygit's own, which is as wide as the whole
		// window rather than the view. With none, it falls back to the COLUMNS it is
		// given.
		t.Views().Main().
			Content(Contains("terminal: none"))
	},
})
