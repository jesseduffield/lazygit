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
		shell.CloneIntoRemote("origin")
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
			Press(keys.Universal.Push).
			Tap(func() {
				t.ExpectPopup().Prompt().
					Title(Equals("Remote to push tag 'new-tag' to:")).
					InitialText(Equals("origin")).
					SuggestionLines(
						Contains("origin"),
					).
					Confirm()
			}).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete tag 'new-tag'?")).
					Select(Contains("Delete remote tag")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().Prompt().
					Title(Equals("Remote from which to remove tag 'new-tag':")).
					InitialText(Equals("origin")).
					SuggestionLines(
						Contains("origin"),
					).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().
					Confirmation().
					Title(Equals("Delete tag 'new-tag'?")).
					Content(Equals("Are you sure you want to delete the remote tag 'new-tag' from 'origin'?")).
					Confirm()
				t.ExpectToast(Equals("Remote tag deleted"))
			}).
			Lines(
				MatchesRegexp(`new-tag.*message`).IsSelected(),
			).
			Tap(func() {
				t.Git().
					RemoteTagDeleted("origin", "new-tag")
			}).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete tag 'new-tag'?")).
					Select(Contains("Delete local tag")).
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
