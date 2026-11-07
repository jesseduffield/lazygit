package am

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ContinueAfterResolvingConflicts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Continue applying patches with git am after resolving the conflicts of one of them",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		setUpConflictingAm(shell, false, true)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("(applying patches) repo → master"))
		t.Views().Information().Content(Contains("Applying patches"))
		t.Views().Options().Content(Contains("View 'am' options: m"))

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU").Contains("file"),
			)

		t.Views().Commits().
			Focus().
			Lines(
				Contains("─── Pending patches"),
				Contains("pick").Contains("add file2"),
				Contains("pick").Contains("<-- CONFLICT --- change file").IsSelected(),
				Contains("─── Commits"),
				Contains("conflicting change"),
				Contains("base"),
			).
			NavigateToLine(Contains("add file2")).
			Tap(func() {
				t.Views().Main().Content(
					Contains("Subject: [PATCH 2/2] add file2").
						Contains("+two"))
			}).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectToast(Equals("Disabled: This action is not allowed while applying patches"))
			})

		t.Shell().UpdateFileAndAdd("file", "resolved\n")

		t.GlobalPress(keys.Universal.CreateRebaseOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Applying patches")).
			Lines(
				Contains("continue"),
				Contains("abort"),
				Contains("skip"),
				Contains("Cancel"),
			).
			Select(Contains("continue")).
			Confirm()

		t.Views().Status().Content(Equals("repo → master"))

		t.Views().Commits().
			Lines(
				Contains("add file2"),
				Contains("change file"),
				Contains("conflicting change"),
				Contains("base"),
			)
	},
})
