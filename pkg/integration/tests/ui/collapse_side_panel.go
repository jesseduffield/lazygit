package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CollapseSidePanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Alt+number collapses a side panel; pressing again uncollapses it",
	SetupConfig: func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(3)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Start with files focused (default)
		t.Views().Files().IsFocused()

		// Alt+2 to collapse the files panel; focus should move away
		t.GlobalPress("<a-2>")
		t.Views().Status().IsFocused()

		// Alt+2 again to uncollapse files and focus it
		t.GlobalPress("<a-2>")
		t.Views().Files().IsFocused()

		// Navigate to branches panel
		t.Views().Branches().Focus()

		// Alt+3 to collapse branches; focus should move away
		t.GlobalPress("<a-3>")
		t.Views().Status().IsFocused()

		// Alt+3 again to uncollapse branches and focus it
		t.GlobalPress("<a-3>")
		t.Views().Branches().IsFocused()

		// Navigate to commits panel
		t.Views().Commits().Focus()

		// Alt+4 to collapse commits; focus should move away
		t.GlobalPress("<a-4>")
		t.Views().Status().IsFocused()

		// When files is collapsed and we jump to it via number key, it uncollapses
		t.Views().Files().Focus()
		t.Views().Files().IsFocused()
	},
})
