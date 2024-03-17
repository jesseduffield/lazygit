package archive

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var Branch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create an archive of a branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("first-change-branch").IsSelected(),
				Contains("second-change-branch"),
				Contains("original-branch"),
			).
			Press(keys.Branches.Archive)

		t.ExpectPopup().Prompt().
			Title(Equals("Choose name for archive of first-change-branch")).
			Type("test-archive-v1.2.3").
			Confirm()

		t.ExpectPopup().Menu().
			Title(Equals("Select archive format")).
			Select(Contains("zip")).
			Confirm()

		t.FileSystem().PathPresent("test-archive-v1.2.3.zip")
	},
})
