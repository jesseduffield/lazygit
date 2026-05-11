package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectedSubmodule = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Use the {{ .SelectedSubmodule }} template variable",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("Initial commit")
		shell.CloneIntoSubmodule("submodule", "path/submodule")
		shell.Commit("Add submodule")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "X",
				Context: "submodules",
				Command: "printf '%s' '{{ .SelectedSubmodule.Path }}' > file.txt",
			},
			{
				Key:     "U",
				Context: "submodules",
				Command: "printf '%s' '{{ .SelectedSubmodule.Url }}' > file.txt",
			},
			{
				Key:     "N",
				Context: "submodules",
				Command: "printf '%s' '{{ .SelectedSubmodule.Name }}' > file.txt",
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Submodules().
			Focus().
			Lines(
				Contains("submodule").IsSelected(),
			)

		t.Views().Submodules().Press("X")
		t.FileSystem().FileContent("file.txt", Equals("path/submodule"))

		t.Views().Submodules().Press("U")
		t.FileSystem().FileContent("file.txt", Equals("../submodule"))

		t.Views().Submodules().Press("N")
		t.FileSystem().FileContent("file.txt", Equals("submodule"))
	},
})
