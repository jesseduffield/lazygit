package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RenamedFilePartial = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select part of a renamed file's changes into a custom patch and remove it from the commit, keeping the rename in place",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
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
			PressEnter()

		// The main view shows the rename together with its content change.
		t.Views().PatchBuilding().
			IsFocused().
			Content(Contains("rename from original").Contains("rename to renamed")).
			ContainsLines(
				Contains(" line1"),
				Contains("-line2"),
				Contains("+line2 changed"),
				Contains(" line3"),
			).
			// Add the hunk (a line selection, as opposed to adding the whole
			// file), so this is a partial patch.
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Common().SelectPatchOption(Contains("Remove patch from original commit"))

		// The rename is preserved; only the content change is gone, so the file
		// is still shown as a rename but now has no content change.
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("original → renamed").IsSelected(),
			)

		t.Views().Main().
			Content(DoesNotContain("line2 changed"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("rename with modification").IsSelected(),
				Contains("first commit"),
			)
	},
})
