package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var ModeSpecificKeybindingSuggestions = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When in various modes, we should corresponding keybinding suggestions onscreen",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(2)
		shell.NewBranch("base-branch")
		shared.MergeConflictsSetup(shell)
		shell.Checkout("base-branch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		rebaseSuggestion := "View rebase options: m"
		cherryPickSuggestion := "Paste (cherry-pick): V"
		bisectSuggestion := "View bisect options: b"
		customPatchSuggestion := "View custom patch options: <c-p>"
		mergeSuggestion := "View merge options: m"

		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 02").IsSelected(),
				Contains("commit 01"),
			).
			Tap(func() {
				// These suggestions are mode-specific so are not shown by default
				t.Views().Options().Content(
					DoesNotContain(rebaseSuggestion).
						DoesNotContain(mergeSuggestion).
						DoesNotContain(cherryPickSuggestion).
						DoesNotContain(bisectSuggestion).
						DoesNotContain(customPatchSuggestion),
				)
			}).
			// Start an interactive rebase
			Press(keys.Universal.Edit).
			Tap(func() {
				// Confirm the rebase suggestion now appears
				t.Views().Options().Content(Contains(rebaseSuggestion))
			}).
			Press(keys.Commits.CherryPickCopy).
			Tap(func() {
				// Confirm the cherry pick suggestion now appears
				t.Views().Options().Content(Contains(cherryPickSuggestion))
				// Importantly, we show multiple of these suggestions at once
				t.Views().Options().Content(Contains(rebaseSuggestion))
			}).
			// Cancel the cherry pick
			PressEscape().
			Tap(func() {
				t.Views().Options().Content(DoesNotContain(cherryPickSuggestion))
			}).
			// Cancel the rebase
			Tap(func() {
				t.Common().AbortRebase()

				t.Views().Options().Content(DoesNotContain(rebaseSuggestion))
			}).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Bisect")).
					Select(MatchesRegexp("Mark.* as bad")).
					Confirm()

				t.Views().Options().Content(Contains(bisectSuggestion))

				// Cancel bisect
				t.Common().ResetBisect()

				t.Views().Options().Content(DoesNotContain(bisectSuggestion))
			}).
			// Enter commit files view
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			// Add a commit file to the patch
			Press(keys.Universal.Select).
			Tap(func() {
				t.Views().Options().Content(Contains(customPatchSuggestion))

				t.Common().ResetCustomPatch()

				t.Views().Options().Content(DoesNotContain(customPatchSuggestion))
			})

		// Test merge options  suggestion
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("first-change-branch")).
			Press(keys.Universal.Select).
			NavigateToLine(Contains("second-change-branch")).
			Press(keys.Branches.MergeIntoCurrentBranch).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Merge")).
					Content(Contains("Are you sure you want to merge")).
					Confirm()

				t.Common().AcknowledgeConflicts()

				t.Views().Options().Content(Contains(mergeSuggestion))

				t.Common().AbortMerge()

				t.Views().Options().Content(DoesNotContain(mergeSuggestion))
			})
	},
})
