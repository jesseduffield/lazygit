package archive

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var Commit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create an archive of a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("original")).
			Press(keys.Branches.Archive)

		t.ExpectPopup().Prompt().
			Title(Contains("Choose name for archive of")).
			Type("test-archive-from-commit").
			Confirm()

		t.ExpectPopup().Menu().
			Title(Equals("Select archive format")).
			Select(Contains("tgz")).
			Confirm()

		t.FileSystem().PathPresent("test-archive-from-commit.tgz")
	},
})
