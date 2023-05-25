package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardStagedChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding staged changes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("fileToRemove", "original content")
		shell.CreateFileAndAdd("file2", "original content")
		shell.Commit("first commit")

		shell.CreateFile("file3", "original content")
		shell.UpdateFile("fileToRemove", "new content")
		shell.UpdateFile("file2", "new content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains(` M file2`).IsSelected(),
				Contains(`?? file3`),
				Contains(` M fileToRemove`),
			).
			NavigateToLine(Contains(`fileToRemove`)).
			PressPrimaryAction().
			Lines(
				Contains(` M file2`),
				Contains(`?? file3`),
				Contains(`M  fileToRemove`).IsSelected(),
			).
			Press(keys.Files.ViewResetOptions)

		t.ExpectPopup().Menu().Title(Equals("")).Select(Contains("Discard staged changes")).Confirm()

		// staged file has been removed
		t.Views().Files().
			Lines(
				Contains(` M file2`),
				Contains(`?? file3`).IsSelected(),
			)

		// the file should have the same content that it originally had, given that that was committed already
		t.FileSystem().FileContent("fileToRemove", Equals("original content"))
	},
})
