package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageAllWithDirtySubmodule = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A submodule with only dirty content (which can't be staged) must not break the stage-all toggle: pressing it repeatedly should keep toggling the other files between staged and unstaged.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("my_submodule_name", "my_submodule_path")
		shell.GitAddAll()
		shell.Commit("add submodule")

		// A submodule with dirty content but no new commit (can't be staged),
		// alongside a regular file that can.
		shell.CreateFile("my_submodule_path/dirty_file", "dirty content")
		shell.CreateFile("regular_file", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus().
			Lines(
				Equals(" M my_submodule_path (submodule)"),
				Equals("?? regular_file"),
			).
			// Stage all: the regular file gets staged; the submodule can't be.
			Press(keys.Files.ToggleStagedAll).
			Lines(
				Equals(" M my_submodule_path (submodule)"),
				Equals("A  regular_file"),
			).
			// Stage all again: nothing is stageable, but the regular file is
			// staged, so this unstages it rather than erroring on the submodule.
			Press(keys.Files.ToggleStagedAll).
			Lines(
				Equals(" M my_submodule_path (submodule)"),
				Equals("?? regular_file"),
			)
	},
})
