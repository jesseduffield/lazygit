package tag

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CrudLightweight = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create and delete a lightweight tag in the tags panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Tags().
			Focus().
			IsEmpty().
			Press(keys.Universal.New).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Create tag")).
					Select(Contains("Lightweight")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("Tag name:")).
					Type("new-tag").
					Confirm()
			}).
			Lines(
				MatchesRegexp(`new-tag.*initial commit`).IsSelected(),
			).
			PressEnter().
			Tap(func() {
				// view the commits of the tag
				t.Views().SubCommits().IsFocused().
					Lines(
						Contains("initial commit"),
					).
					PressEscape()
			}).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Delete tag")).
					Content(Equals("Are you sure you want to delete tag 'new-tag'?")).
					Confirm()
			}).
			IsEmpty()
	},
})
