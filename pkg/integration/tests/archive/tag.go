package archive

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var Tag = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create an archive of a tag",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
		shell.CreateLightweightTag("tag", "HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Tags().
			Focus().
			Lines(
				Contains("tag").IsSelected(),
			).Press(keys.Branches.Archive)

		t.ExpectPopup().Prompt().
			Title(Equals("Choose name for archive of tag")).
			Type("test-archive-tag").
			Confirm()

		t.ExpectPopup().Menu().
			Title(Equals("Select archive format")).
			Select(Contains("zip")).
			Confirm()

		t.FileSystem().PathPresent("test-archive-tag.zip")
	},
})
