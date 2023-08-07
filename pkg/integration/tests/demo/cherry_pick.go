package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPick = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.Gui.NerdFontsVersion = "3"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommitsWithRandomMessages(50)

		shell.
			EmptyCommit("Fix bug in timezone conversion.").
			NewBranch("hotfix/fix-bug").
			NewBranch("feature/user-module").
			Checkout("hotfix/fix-bug").
			EmptyCommit("Integrate support for markdown in user posts").
			EmptyCommit("Remove unused code and libraries").
			Checkout("feature/user-module").
			EmptyCommit("Handle session timeout gracefully").
			EmptyCommit("Add Webpack for asset bundling").
			Checkout("hotfix/fix-bug")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Cherry pick commits from another branch")
		t.Wait(1000)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("hotfix/fix-bug"),
				Contains("feature/user-module"),
				Contains("master"),
			).
			SelectNextItem().
			Wait(300).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			TopLines(
				Contains("Add Webpack for asset bundling").IsSelected(),
				Contains("Handle session timeout gracefully"),
				Contains("Fix bug in timezone conversion."),
			).
			Press(keys.Commits.CherryPickCopy).
			Tap(func() {
				t.Views().Information().Content(Contains("1 commit copied"))
			}).
			SelectNextItem().
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("2 commits copied"))

		t.Views().Commits().
			Focus().
			TopLines(
				Contains("Remove unused code and libraries").IsSelected(),
				Contains("Integrate support for markdown in user posts"),
				Contains("Fix bug in timezone conversion."),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.Wait(1000)
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Contains("Are you sure you want to cherry-pick the copied commits onto this branch?")).
					Confirm()
			}).
			TopLines(
				Contains("Add Webpack for asset bundling"),
				Contains("Handle session timeout gracefully"),
				Contains("Remove unused code and libraries"),
				Contains("Integrate support for markdown in user posts"),
				Contains("Fix bug in timezone conversion."),
			).
			Tap(func() {
				// we need to manually exit out of cherry pick mode
				t.Views().Information().Content(Contains("2 commits copied"))
			}).
			PressEscape().
			Tap(func() {
				t.Views().Information().Content(DoesNotContain("commits copied"))
			})
	},
})
