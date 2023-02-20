package tag

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CrudAnnotated = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create and delete an annotated tag in the tags panel",
	ExtraCmdArgs: "",
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
					Select(Contains("annotated")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("Tag name:")).
					Type("new-tag").
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("Tag message:")).
					Type("message").
					Confirm()
			}).
			Lines(
				MatchesRegexp(`new-tag.*message`).IsSelected(),
			).
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
