package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RenamedFileWhole = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add a whole renamed file to a custom patch and remove it from the commit, taking the rename with it",
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
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		// The whole file is added, so the patch carries the rename itself.
		t.Views().Secondary().
			ContainsLines(
				Contains("rename from original"),
				Contains("rename to renamed"),
			)

		t.Common().SelectPatchOption(Contains("Remove patch from original commit"))

		// The rename went with the patch, so the commit no longer touches the file.
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("(none)"),
			)

		t.Views().Commits().
			Focus().
			Lines(
				Contains("rename with modification").IsSelected(),
				Contains("first commit"),
			)
	},
})
