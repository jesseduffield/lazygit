package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EnterNested = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Enter a nested submodule",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		setupNestedSubmodules(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Submodules().Focus().
			Lines(
				Equals("outerSubName").IsSelected(),
				Equals("  - innerSubName"),
			).
			Tap(func() {
				t.Views().Main().ContainsLines(
					Contains("Name: outerSubName"),
					Contains("Path: modules/outerSubPath"),
					Contains("Url:  ../outerSubmodule"),
				)
			}).
			SelectNextItem().
			Tap(func() {
				t.Views().Main().ContainsLines(
					Contains("Name: outerSubName/innerSubName"),
					Contains("Path: modules/outerSubPath/modules/innerSubPath"),
					Contains("Url:  ../innerSubmodule"),
				)
			}).
			// enter the nested submodule
			PressEnter()

		if t.Git().Version().IsAtLeast(2, 22, 0) {
			t.Views().Status().Content(Contains("innerSubPath(innerSubName)"))
		} else {
			t.Views().Status().Content(Contains("innerSubPath"))
		}
		t.Views().Commits().ContainsLines(
			Contains("initial inner commit"),
		)

		t.Views().Files().PressEscape()
		t.Views().Status().Content(Contains("repo"))
	},
})
