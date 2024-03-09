package submodule

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoveNested = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a nested submodule",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		setupNestedSubmodules(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		gitDirSubmodulePath, _ := filepath.Abs(".git/modules/outerSubName/modules/innerSubName")
		t.FileSystem().PathPresent(gitDirSubmodulePath)

		t.Views().Submodules().Focus().
			Lines(
				Equals("outerSubName").IsSelected(),
				Equals("  - innerSubName"),
			).
			SelectNextItem().
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Remove submodule")).
					Content(Equals("Are you sure you want to remove submodule 'outerSubName/innerSubName' and its corresponding directory? This is irreversible.")).
					Confirm()
			}).
			Lines(
				Equals("outerSubName").IsSelected(),
			).
			Press(keys.Universal.GoInto)

		t.Views().Files().IsFocused().
			Lines(
				Contains("modules").IsSelected(),
				MatchesRegexp(`D.*innerSubPath`),
				MatchesRegexp(`M.*\.gitmodules`),
			).
			NavigateToLine(Contains(".gitmodules"))

		t.Views().Main().Content(
			Contains("-[submodule \"innerSubName\"]").
				Contains("-   path = modules/innerSubPath").
				Contains("-   url = ../innerSubmodule"),
		)

		t.FileSystem().PathNotPresent(gitDirSubmodulePath)
	},
})
