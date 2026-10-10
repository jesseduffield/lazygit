package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GlobalContext = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ensure global context works",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("my change")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:         config.Keybinding{"X"},
				Context:     "global",
				Command:     "touch myfile",
				Description: "Global custom command",
				DisplayHint: true,
			},
			{
				Key:         config.Keybinding{"e"},
				Context:     "global",
				Command:     "touch shadowedfile",
				Description: "Shadowed custom command",
				DisplayHint: true,
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		hint := "Global custom command"

		// commits
		t.Views().Commits().
			Focus().
			Tap(func() {
				t.Views().Options().Content(Contains(hint))
			}).
			Press(config.Keybinding{"X"})

		t.Views().Files().
			Focus().
			Tap(func() {
				t.Views().Options().Content(Contains(hint))
			}).
			Lines(Contains("myfile"))

		t.Shell().DeleteFile("myfile")
		t.GlobalPress(keys.Files.RefreshFiles)

		// branches
		t.Views().Branches().
			Focus().
			Tap(func() {
				t.Views().Options().Content(Contains(hint))
			}).
			Press(config.Keybinding{"X"})

		t.Views().Files().
			Focus().
			Lines(Contains("myfile"))

		t.Shell().DeleteFile("myfile")
		t.GlobalPress(keys.Files.RefreshFiles)

		// files
		t.Views().Files().
			Focus().
			Tap(func() {
				t.Views().Options().Content(Contains(hint))
				t.Views().Options().Content(DoesNotContain("Shadowed custom command"))
			}).
			Press(config.Keybinding{"X"})

		t.Views().Files().
			Focus().
			Lines(Contains("myfile"))

		t.Shell().DeleteFile("myfile")

		// search
		t.Views().Commits().
			Focus().
			Press(keys.Universal.StartSearch).
			Tap(func() {
				t.Views().Options().Content(DoesNotContain(hint))
			})
	},
})
