package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Remove = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a submodule",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("my_submodule")
		shell.GitAddAll()
		shell.Commit("add submodule")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Submodules().Focus().
			Lines(
				Contains("my_submodule").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Remove submodule")).
					Content(Equals("Are you sure you want to remove submodule 'my_submodule' and its corresponding directory? This is irreversible.")).
					Confirm()
			}).
			IsEmpty()

		t.Views().Files().Focus().
			Lines(
				MatchesRegexp(`M.*\.gitmodules`).IsSelected(),
				MatchesRegexp(`D.*my_submodule`),
			)

		t.Views().Main().Content(
			Contains("-[submodule \"my_submodule\"]").
				Contains("-   path = my_submodule").
				Contains("-   url = ../other_repo"),
		)
	},
})
