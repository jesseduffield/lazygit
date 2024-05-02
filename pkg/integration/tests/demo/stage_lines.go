package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var originalFile = `# Lazygit

Simple terminal UI for git commands

![demo](https://user-images.gh.com/demo.gif)

## Installation

### Homebrew

`

var updatedFile = `# Lazygit

Simple terminal UI for git
(Not too simple though)

![demo](https://user-images.gh.com/demo.gif)

## Installation

### Homebrew

Just do brew install lazygit and bada bing bada
boom you have begun on the path of laziness.

`

var StageLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage individual lines",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		setDefaultDemoConfig(config)
		config.UserConfig.Gui.ShowFileTree = false
		config.UserConfig.Gui.ShowCommandLog = false
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("docs-fix")
		shell.CreateNCommitsWithRandomMessages(30)
		shell.CreateFileAndAdd("docs/README.md", originalFile)
		shell.Commit("Update docs/README")
		shell.UpdateFile("docs/README.md", updatedFile)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Stage individual lines")
		t.Wait(1000)

		t.Views().Files().
			IsFocused().
			PressEnter()

		t.Views().Staging().
			IsFocused().
			Press(keys.Universal.ToggleRangeSelect).
			PressFast(keys.Universal.NextItem).
			PressFast(keys.Universal.NextItem).
			Wait(500).
			PressPrimaryAction().
			Wait(500).
			PressEscape()

		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChanges).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Type("Update tagline").
					Confirm()
			})

		t.Views().Commits().
			Focus()
	},
})
