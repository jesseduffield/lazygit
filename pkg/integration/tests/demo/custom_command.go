package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var customCommandContent = `
customCommands:
  - key: 'a'
    command: 'git checkout {{.Form.Branch}}'
    context: 'localBranches'
    prompts:
    - type: 'input'
      title: 'Enter a branch name to checkout:'
      key: 'Branch'
			suggestions:
				preset: 'branches'
`

var CustomCommand = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Invoke a custom command",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(cfg *config.AppConfig) {
		// No idea why I had to use version 2: it should be using my own computer's
		// font and the one iterm uses is version 3.
		cfg.UserConfig.Gui.NerdFontsVersion = "2"

		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "localBranches",
				Command: `git checkout {{.Form.Branch}}`,
				Prompts: []config.CustomCommandPrompt{
					{
						Key:   "Branch",
						Type:  "input",
						Title: "Enter a branch name to checkout",
						Suggestions: config.CustomCommandSuggestions{
							Preset: "branches",
						},
					},
				},
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommitsWithRandomMessages(30)
		shell.NewBranch("feature/user-authentication")
		shell.NewBranch("feature/payment-processing")
		shell.NewBranch("feature/search-functionality")
		shell.NewBranch("feature/mobile-responsive")
		shell.EmptyCommit("Make mobile response")
		shell.NewBranch("bugfix/fix-login-issue")
		shell.HardReset("HEAD~1")
		shell.NewBranch("bugfix/fix-crash-bug")
		shell.CreateFile("custom_commands_example.yml", customCommandContent)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Invoke a custom command")
		t.Wait(1500)

		t.Views().Branches().
			Focus().
			Wait(500).
			Press("a").
			Tap(func() {
				t.Wait(500)

				t.ExpectPopup().Prompt().
					Title(Equals("Enter a branch name to checkout")).
					Type("mobile").
					ConfirmFirstSuggestion()
			})
	},
})
