package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CycleDiffRenderers = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cycle forwards and backwards through configured diff renderers",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			// an explicit name overrides the derived one
			{Name: "custom name", Command: "cat"},
			// no name, so it's derived from the first word of the command
			{Command: "cat -n"},
			// rawGit derives it from the first argument if any
			{Type: "rawGit", Args: []string{"--color-words"}},
			// neither name nor command, so it falls back to the default label
			{Type: "rawGit"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(1)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.CycleDiffRenderers)
		t.ExpectToast(Equals("Diff renderer: cat (2 of 4)"))

		t.Views().Commits().Press(keys.Universal.CycleDiffRenderers)
		t.ExpectToast(Equals("Diff renderer: --color-words (3 of 4)"))

		t.Views().Commits().Press(keys.Universal.CycleDiffRenderers)
		t.ExpectToast(Equals("Diff renderer: (default) (4 of 4)"))

		// cycling forward past the last diff renderer wraps around to the first
		t.Views().Commits().Press(keys.Universal.CycleDiffRenderers)
		t.ExpectToast(Equals("Diff renderer: custom name (1 of 4)"))

		// cycling backward past the first diff renderer wraps around to the last
		t.Views().Commits().Press(keys.Universal.CycleDiffRenderersReverse)
		t.ExpectToast(Equals("Diff renderer: (default) (4 of 4)"))

		t.Views().Commits().Press(keys.Universal.CycleDiffRenderersReverse)
		t.ExpectToast(Equals("Diff renderer: --color-words (3 of 4)"))
	},
})
