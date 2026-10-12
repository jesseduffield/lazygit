package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardRemovedSubmodule = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard the directory that a removed submodule leaves behind",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("sub", "modules/sub")
		shell.Commit("add submodule")

		// This is what checking out a commit that removes the submodule does:
		// the submodule is gone from the index and from .gitmodules, but git
		// leaves its directory behind.
		shell.RunCommand([]string{"git", "rm", "--cached", "modules/sub"})
		shell.RunCommand([]string{"git", "rm", ".gitmodules"})
		shell.Commit("remove submodule")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ modules").IsSelected(),
				Equals("  ?? sub (removed submodule)"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			IsFocused().
			IsEmpty()

		t.FileSystem().PathNotPresent("modules/sub")
	},
})
