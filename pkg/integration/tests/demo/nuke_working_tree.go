package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NukeWorkingTree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Nuke the working tree",
	ExtraCmdArgs: []string{"status"},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		setDefaultDemoConfig(config)
		config.UserConfig.Gui.AnimateExplosion = true
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("blah")
		shell.CreateFile("controllers/red_controller.rb", "")
		shell.CreateFile("controllers/green_controller.rb", "")
		shell.CreateFileAndAdd("controllers/blue_controller.rb", "")
		shell.CreateFile("controllers/README.md", "")
		shell.CreateFileAndAdd("views/helpers/list.rb", "")
		shell.CreateFile("views/helpers/sort.rb", "")
		shell.CreateFileAndAdd("views/users_view.rb", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Nuke the working tree")
		t.Wait(1000)

		t.Views().Files().
			IsFocused().
			Wait(1000).
			Press(keys.Files.ViewResetOptions).
			Tap(func() {
				t.Wait(1000)

				t.ExpectPopup().Menu().
					Title(Equals("")).
					Select(Contains("Nuke working tree")).
					Confirm()
			})
	},
})
