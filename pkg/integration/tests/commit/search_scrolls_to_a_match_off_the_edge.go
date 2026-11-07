package commit

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SearchScrollsToAMatchOffTheEdge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Searching the commits panel scrolls it sideways to show a match past its right edge, also in a commit that was below the visible area",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch")
		for i := 1; i <= 30; i++ {
			message := fmt.Sprintf("commit %02d", i)
			switch i {
			case 29:
				message += " " + strings.Repeat("word ", 20) + "needle"
			case 3:
				message += " " + strings.Repeat("word ", 40) + "needle"
			}
			shell.EmptyCommit(message)
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			FilterOrSearch("needle").
			SelectedLines(
				Contains("commit 29"),
			).
			OriginX(112).
			// Commit 03 is below the visible area, and its match is 100 columns
			// further right than the one in commit 29.
			Press(keys.Universal.NextMatch).
			SelectedLines(
				Contains("commit 03"),
			).
			OriginX(212)
	},
})
