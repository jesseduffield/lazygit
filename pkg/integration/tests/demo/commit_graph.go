package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitGraph = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show commit graph",
	ExtraCmdArgs: []string{"log"},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.Gui.NerdFontsVersion = "3"
		config.UserConfig.Gui.AuthorColors = map[string]string{
			"Fredrica Greenhill": "#fb5aa3",
			"Oscar Reuenthal":    "#86c82f",
			"Paul Oberstein":     "#ffd500",
			"Siegfried Kircheis": "#fe7e11",
			"Yang Wen-li":        "#8e3ccb",
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateRepoHistory()
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("View commit log")
		t.Wait(1000)

		t.Views().Commits().
			IsFocused().
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100).
			SelectNextItem().
			Wait(100)
	},
})
