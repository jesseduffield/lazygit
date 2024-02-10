package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Remove = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a submodule",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("my_submodule_name", "my_submodule_path")
		shell.GitAddAll()
		shell.Commit("add submodule")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		gitDirSubmodulePath := ".git/modules/my_submodule_name"
		t.FileSystem().PathPresent(gitDirSubmodulePath)

		t.Views().Submodules().Focus().
			Lines(
				Contains("my_submodule_name").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Remove submodule")).
					Content(Equals("Are you sure you want to remove submodule 'my_submodule_name' and its corresponding directory? This is irreversible.")).
					Confirm()
			}).
			IsEmpty()

		t.Views().Files().Focus().
			Lines(
				MatchesRegexp(`M.*\.gitmodules`).IsSelected(),
				MatchesRegexp(`D.*my_submodule_path`),
			)

		t.Views().Main().Content(
			Contains("-[submodule \"my_submodule_name\"]").
				Contains("-   path = my_submodule_path").
				Contains("-   url = ../other_repo"),
		)

		t.FileSystem().PathNotPresent(gitDirSubmodulePath)
	},
})
