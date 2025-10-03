package demo

import (
	"os"
	"os/exec"
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPick = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		setDefaultDemoConfig(config)
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
		wd, err := os.Getwd()
		if err != nil {
			t.Fail("Could not determine working directory: " + err.Error())
			return
		}

		cherryPickInProgress := func() bool {
			cmd := exec.Command("git", "rev-parse", "CHERRY_PICK_HEAD")
			cmd.Dir = wd

			return cmd.Run() == nil
		}

		waitForCherryPickInProgress := func(timeout time.Duration) bool {
			deadline := time.Now().Add(timeout)
			for {
				if cherryPickInProgress() {
					return true
				}

				if time.Now().After(deadline) {
					return false
				}

				t.Wait(100)
			}
		}

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
					Content(Contains("Are you sure you want to cherry-pick the 2 copied commit(s) onto this branch?")).
					Confirm()
			}).
			Tap(func() {
				if !waitForCherryPickInProgress(time.Second) {
					return
				}

				for {
					t.ExpectPopup().Menu().
						Title(Equals("Cherry-pick produced no changes")).
						ContainsLines(
							Contains("Skip this cherry-pick"),
							Contains("Create empty commit and continue"),
							Contains("Cancel"),
						).
						Select(Contains("Create empty commit and continue")).
						Confirm()

					t.Wait(100)

					if !waitForCherryPickInProgress(time.Second) {
						break
					}
				}
			}).
			Tap(func() {
				t.Shell().RunCommandExpectError([]string{"git", "rev-parse", "CHERRY_PICK_HEAD"})
			}).
			TopLines(
				Contains("Add Webpack for asset bundling"),
				Contains("Handle session timeout gracefully"),
				Contains("Remove unused code and libraries"),
				Contains("Integrate support for markdown in user posts"),
				Contains("Fix bug in timezone conversion."),
			).
			Tap(func() {
				t.Views().Information().Content(DoesNotContain("commits copied"))
			})
	},
})
