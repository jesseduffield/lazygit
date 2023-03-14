package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Add = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add a submodule",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.Clone("other_repo")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Submodules().Focus().
			Press(keys.Universal.New).
			Tap(func() {
				t.ExpectPopup().Prompt().
					Title(Equals("new submodule URL:")).
					Type("../other_repo").Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("new submodule name:")).
					InitialText(Equals("other_repo")).
					Clear().Type("my_submodule").Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("new submodule path:")).
					InitialText(Equals("my_submodule")).
					Clear().Type("my_submodule_path").Confirm()
			}).
			Lines(
				Contains("my_submodule").IsSelected(),
			)

		t.Views().Main().TopLines(
			Contains("Name: my_submodule"),
			Contains("Path: my_submodule_path"),
			Contains("Url:  ../other_repo"),
		)

		t.Views().Files().Focus().
			Lines(
				Contains(".gitmodules").IsSelected(),
				Contains("my_submodule_path (submodule)"),
			).
			Tap(func() {
				t.Views().Main().Content(
					Contains("[submodule \"my_submodule\"]").
						Contains("path = my_submodule_path").
						Contains("url = ../other_repo"),
				)
			}).
			SelectNextItem().
			Tap(func() {
				t.Views().Main().Content(
					Contains("Submodule my_submodule_path").
						Contains("(new submodule)"),
				)
			})
	},
})
