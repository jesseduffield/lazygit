package conflicts

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var ConflictMarkerSizeResolve = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Resolves a conflict in a file whose conflict-marker-size gitattribute makes its markers longer than usual",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.SetCustomConflictMarkerSize(shell)
		shared.CreateMergeConflictFileMultiple(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		startMarker := strings.Repeat("<", shared.CustomConflictMarkerSize)

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file").IsSelected(),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			SelectedLines(
				Contains(startMarker+" HEAD"),
				Contains("First Change"),
				Contains(strings.Repeat("=", shared.CustomConflictMarkerSize)),
			).
			PressPrimaryAction().
			Content(DoesNotContain(startMarker + " HEAD\nFirst Change"))
	},
})
