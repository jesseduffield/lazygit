package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var IgnoreSuspiciousName = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ignore a submodule whose name in .gitmodules would put its git dir outside the repository",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("../victim/file", "content")

		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("my_submodule_name", "my_submodule_path")
		shell.GitAddAll()
		shell.Commit("add submodule")

		// The git dir of a submodule is .git/modules/<name>, so this name
		// points at the victim directory next to the repository
		shell.UpdateFileAndAdd(".gitmodules", `[submodule "../../../victim"]
	path = my_submodule_path
	url = ../my_submodule_name
`)
		shell.Commit("rename submodule")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Submodules().Focus().
			IsEmpty()

		t.FileSystem().
			PathPresent("../victim/file")
	},
})
