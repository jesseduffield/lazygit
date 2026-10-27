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
TODO: mention the other install methods
`

var StageHunksOrLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage hunks or individual lines",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		setDefaultDemoConfig(config)
		config.GetUserConfig().Gui.ShowFileTree = false
		config.GetUserConfig().Gui.ShowCommandLog = false
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("docs-fix")
		shell.CreateNCommitsWithRandomMessages(30)
		shell.CreateFileAndAdd("docs/README.md", originalFile)
		shell.Commit("Update docs/README")
		shell.UpdateFile("docs/README.md", updatedFile)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Stage a hunk")
		t.Wait(1000)

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-Simple terminal UI for git commands"),
				Contains("+Simple terminal UI for git"),
				Contains("+(Not too simple though)"),
			).
			Wait(1000).
			PressPrimaryAction().
			Wait(1000).
			SetCaptionPrefix("Stage individual lines").
			Press(keys.Main.ToggleSelectHunk).
			Wait(500).
			Press(keys.Universal.ToggleRangeSelect).
			PressFast(keys.Universal.NextItem).
			SelectedLines(
				Contains("+Just do brew install lazygit"),
				Contains("+boom you have begun"),
			).
			Wait(500).
			PressPrimaryAction().
			Wait(1000).
			SelectedLines(
				Contains("+TODO: mention the other install methods"),
			).
			SetCaptionPrefix("Commit our changes").
			Press(keys.Files.CommitChanges).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Type("Update tagline and install instructions").
					Confirm()
			})

		t.Views().Commits().
			Focus()
	},
})
