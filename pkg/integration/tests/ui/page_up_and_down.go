package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

const (
	// The height of the commits panel in this test's window, in lines.
	commitsPanelHeight = 5
	// Paging keeps one line of overlap between the old and the new page.
	pageDelta = commitsPanelHeight - 1
)

var PageUpAndDown = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Paging down and up keeps the selection at the edge of the viewport",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(40)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLineIdx(0).
			OriginY(0).
			Press(keys.Universal.NextPage).
			// The selection moves to the bottom of the viewport; nothing scrolls yet
			SelectedLineIdx(commitsPanelHeight - 1).
			OriginY(0).
			Press(keys.Universal.NextPage).
			// Now the view scrolls by a page, and the selection stays at the bottom
			SelectedLineIdx(commitsPanelHeight - 1 + pageDelta).
			OriginY(pageDelta).
			Press(keys.Universal.PrevPage).
			// The selection moves to the top of the viewport; nothing scrolls
			SelectedLineIdx(pageDelta).
			OriginY(pageDelta).
			Press(keys.Universal.PrevPage).
			// And back a page, with the selection staying at the top
			SelectedLineIdx(0).
			OriginY(0)
	},
})
