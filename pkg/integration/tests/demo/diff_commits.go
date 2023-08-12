package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiffCommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Diff two commits",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		setDefaultDemoConfig(config)

		config.UserConfig.Gui.ShowFileTree = false
		config.UserConfig.Gui.ShowCommandLog = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommitsWithRandomMessages(50)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Compare two commits")
		t.Wait(1000)

		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("Replace deprecated lifecycle methods in React components")).
			Wait(1000).
			Press(keys.Universal.DiffingMenu).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Diffing")).
					TopLines(
						MatchesRegexp(`Diff .*`),
					).
					Wait(500).
					Confirm()
			}).
			NavigateToLine(Contains("Move constants to a separate config file")).
			Wait(1000).
			PressEnter()
	},
})
