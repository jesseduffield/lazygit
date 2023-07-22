package tag

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CrudAnnotated = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create and delete an annotated tag in the tags panel",
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
				t.ExpectPopup().CommitMessagePanel().
					Title(Equals("Tag name")).
					Type("new-tag").
					SwitchToDescription().
					Title(Equals("Tag description")).
					Type("message").
					SwitchToSummary().
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
			IsEmpty().
			Press(keys.Universal.New).
			Tap(func() {
				// confirm content is cleared on next tag create
				t.ExpectPopup().CommitMessagePanel().
					Title(Equals("Tag name")).
					InitialText(Equals(""))
			})
	},
})
