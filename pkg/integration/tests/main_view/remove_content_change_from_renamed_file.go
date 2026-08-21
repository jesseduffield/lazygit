package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoveContentChangeFromRenamedFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a renamed file's content change from a commit while keeping its rename",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("original", "line1\nline2\nline3\nline4\nline5\n")
		shell.Commit("first commit")

		shell.RenameFileInGit("original", "renamed")
		shell.UpdateFileAndAdd("renamed", "line1\nline2 changed\nline3\nline4\nline5\n")
		shell.Commit("rename with modification")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("rename with modification").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("original → renamed").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Content(Contains("rename from original").Contains("rename to renamed")).
			SelectedLines(Contains("-line2")).
			Press(keys.Universal.ToggleRangeSelect).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("-line2"),
				Contains("+line2 changed"),
			).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		// Only part of the file is in the patch, so the patch leaves the rename behind in
		// the commit and carries the content change alone — which is how the pane beside
		// the diff shows it, over the name the file was renamed to.
		t.Views().Secondary().
			ContainsLines(
				Contains("diff --git a/renamed b/renamed"),
				Contains("index"),
				Contains("--- a/renamed"),
				Contains("+++ b/renamed"),
			).
			ContainsLines(
				Contains(" line1"),
				Contains("-line2"),
				Contains("+line2 changed"),
				Contains(" line3"),
			)

		t.Common().SelectPatchOption(Contains("Remove patch from original commit"))

		t.Views().CommitFiles().Lines(
			Contains("original → renamed").IsSelected(),
		)
		t.Views().Main().
			IsFocused().
			Content(DoesNotContain("line2 changed"))
		t.Views().Commits().Lines(
			Contains("rename with modification").IsSelected(),
			Contains("first commit"),
		)
	},
})
