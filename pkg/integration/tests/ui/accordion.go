package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// When in accordion mode, Lazygit looks like this:
//
// ╶─Status─────────────────────────╴┌─Patch──────────────────────────────────────────────────────────┐
// ╶─Files - Submodules──────0 of 0─╴│commit 6e56dd04b70e548976f7f2928c4d9c359574e2bc                 ▲
// ╶─Local branches - Remotes1 of 1─╴│Author: CI <CI@example.com>                                     █
// ┌─Commits - Reflog───────────────┐│Date:   Wed Jul 19 22:00:03 2023 +1000                          │
// │7fe02805 CI commit 12           ▲│                                                                ▼
// │6e56dd04 CI commit 11           █└────────────────────────────────────────────────────────────────┘
// │a35c687d CI commit 10           ▼┌─Command log────────────────────────────────────────────────────┐
// └───────────────────────10 of 20─┘│Random tip: To filter commits by path, press '<ctrl+s>'         │
// ╶─Stash───────────────────0 of 0─╴└────────────────────────────────────────────────────────────────┘
//  <pgup>/<pgdown>: Scroll, <esc>: Cancel, q: Quit, ?: Keybindings, 1-Donate Ask Question unversioned

var Accordion = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify accordion mode kicks in when the screen height is too small",
	ExtraCmdArgs: []string{},
	Width:        100,
	Height:       10,
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(20)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			VisibleLines(
				Contains("commit 20").IsSelected(),
				Contains("commit 19"),
				Contains("commit 18"),
			).
			// go past commit 11, then come back, so that it ends up in the centre of the viewport
			NavigateToLine(Contains("commit 11")).
			NavigateToLine(Contains("commit 10")).
			NavigateToLine(Contains("commit 11")).
			VisibleLines(
				Contains("commit 12"),
				Contains("commit 11").IsSelected(),
				Contains("commit 10"),
			)

		t.Views().Files().
			Focus()

		// ensure we retain the same viewport upon re-focus
		t.Views().Commits().
			Focus().
			VisibleLines(
				Contains("commit 12"),
				Contains("commit 11").IsSelected(),
				Contains("commit 10"),
			)
	},
})
