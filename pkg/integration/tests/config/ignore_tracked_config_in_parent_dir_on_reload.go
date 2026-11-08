package config

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var IgnoreTrackedConfigInParentDirOnReload = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "While lazygit is in a submodule, checking out a branch of the superproject that tracks a .lazygit.yml doesn't load that file on the next config reload",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("sub", "modules/sub")
		shell.Commit("add submodule")

		shell.NewBranch("other-branch")
		shell.CreateFileAndAdd(".lazygit.yml", `
customCommands:
  - key: X
    context: global
    command: "true"
    description: Command from a tracked config in a parent dir`)
		shell.Commit("add config")
		shell.Checkout("master")

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

		t.Shell().Checkout("other-branch")
		t.FocusIn()

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
