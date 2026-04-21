package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CyclePagers = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cycle forwards and backwards through configured pagers",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Pagers = []config.PagingConfig{
			// an explicit name overrides the derived one
			{Name: "custom name", Pager: "cat"},
			// no name, so it's derived from the first word of the command
			{Pager: "cat -n"},
			// neither name nor command, so it falls back to the default label
			{},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(1)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.CyclePagers)
		t.ExpectToast(Equals("Pager: cat (2 of 3)"))

		t.Views().Commits().Press(keys.Universal.CyclePagers)
		t.ExpectToast(Equals("Pager: (default) (3 of 3)"))

		// cycling forward past the last pager wraps around to the first
		t.Views().Commits().Press(keys.Universal.CyclePagers)
		t.ExpectToast(Equals("Pager: custom name (1 of 3)"))

		// cycling backward past the first pager wraps around to the last
		t.Views().Commits().Press(keys.Universal.CyclePagersReverse)
		t.ExpectToast(Equals("Pager: (default) (3 of 3)"))

		t.Views().Commits().Press(keys.Universal.CyclePagersReverse)
		t.ExpectToast(Equals("Pager: cat (2 of 3)"))
	},
})
