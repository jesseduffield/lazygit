package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var ConflictMarkerSizeNotAutoStaged = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Doesn't auto-stage an unresolved file whose conflict-marker-size gitattribute makes its markers longer than usual",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.SetCustomConflictMarkerSize(shell)
		shared.CreateMergeConflictFile(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Common().PretendMergeOrRebaseStartedInLazygit()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file").IsSelected(),
			).
			// Each refresh checks whether the conflicts are still there
			Press(keys.Universal.Refresh).
			// They are, so the file doesn't get staged and we don't get asked to
			// continue the merge
			Lines(
				Contains("UU file").IsSelected(),
			).
			// Once they really are resolved, we do
			Tap(func() {
				t.Shell().UpdateFile("file", "resolved content")
			}).
			Press(keys.Universal.Refresh).
			Tap(func() {
				t.Common().ContinueOnConflictsResolved("merge")
			}).
			IsEmpty()
	},
})
