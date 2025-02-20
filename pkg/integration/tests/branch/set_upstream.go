package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SetUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Set the upstream of a branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CloneIntoRemote("origin")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Universal.NextScreenMode). // we need to enlargen the window to see the upstream
			Lines(
				Contains("master").DoesNotContain("origin master").IsSelected(),
			).
			Press(keys.Branches.SetUpstream).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Upstream options")).
					Select(Contains(" Set upstream of selected branch")). // using leading space to disambiguate from the 'reset' option
					Confirm()

				// Origin is preselected on the users behalf because there is only one remote

				t.ExpectPopup().Prompt().
					Title(Equals("Enter upstream branch on remote 'origin'")).
					SuggestionLines(Equals("master")).
					ConfirmFirstSuggestion()
			}).
			Lines(
				Contains("master").Contains("origin master").IsSelected(),
			)
	},
})
