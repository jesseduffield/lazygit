package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitMessagePanelOutput = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using custom command output to prefill the commit message panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("myfile", "myfile content")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:         config.Keybinding{"a"},
				Context:     "files",
				Command:     `printf "generated subject\n\ngenerated body\n"`,
				LoadingText: "Generating commit message",
				Output:      "commitMessagePanel",
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(config.Keybinding{"a"})

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			Content(Equals("generated subject")).
			SwitchToDescription().
			Title(Equals("Commit description")).
			Content(Equals("generated body")).
			SwitchToSummary().
			Confirm()

		t.Views().Files().
			IsEmpty()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("generated subject").IsSelected(),
			)

		t.Views().Main().Content(MatchesRegexp("generated subject\n\\s*\n\\s*generated body"))
	},
})
