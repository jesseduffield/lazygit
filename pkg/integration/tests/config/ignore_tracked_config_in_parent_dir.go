package config

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var IgnoreTrackedConfigInParentDir = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When entering a submodule, a .lazygit.yml in a parent directory is ignored if the superproject tracks it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("sub", "modules/sub")
		shell.Commit("add submodule")

		shell.CreateFileAndAdd(".lazygit.yml", `
customCommands:
  - key: X
    context: global
    command: "true"
    description: Command from a tracked config in a parent dir`)
		shell.Commit("add config")

		shell.CreateFile("modules/.lazygit.yml", `
customCommands:
  - key: Y
    context: global
    command: "true"
    description: Command from an untracked config in a parent dir`)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Submodules().Focus().
			Lines(
				Contains("sub").IsSelected(),
			).
			PressEnter()

		t.Views().Status().Content(Contains("sub"))

		t.ExpectToast(MatchesRegexp(`^Ignored config tracked by git: .*repo.\.lazygit\.yml$`))

		t.Views().Files().IsFocused().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Filter("in a parent dir").
			Lines(
				Contains("─── Global"),
				Contains("Command from an untracked config in a parent dir").IsSelected(),
			)
	},
})
