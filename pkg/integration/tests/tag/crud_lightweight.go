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
				MatchesRegexp(`new-tag.*initial commit`).IsSelected(),
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
			IsEmpty()
	},
})
