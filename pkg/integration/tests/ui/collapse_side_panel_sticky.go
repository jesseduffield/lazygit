package ui

import (
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CollapseSidePanelSticky = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Collapsing a side panel is persisted to state.yml so it survives restarts",
	SetupConfig: func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(3)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		stateFile := filepath.Join(os.Getenv("CONFIG_DIR"), "state.yml")

		// Collapse files (pane 2) via alt+2
		t.Views().Files().IsFocused()
		t.GlobalPress("<a-2>")
		t.Views().Status().IsFocused()

		// State file should now record "files" as collapsed
		t.FileSystem().FileContent(stateFile, Contains("collapsedsidewindows"))
		t.FileSystem().FileContent(stateFile, Contains("- files"))

		// Collapse branches (pane 3) too via alt+3
		t.Views().Branches().Focus()
		t.GlobalPress("<a-3>")
		t.Views().Status().IsFocused()

		t.FileSystem().FileContent(stateFile, Contains("- branches"))

		// Uncollapsing should remove the entry from state.yml
		t.GlobalPress("<a-2>")
		t.Views().Files().IsFocused()

		t.FileSystem().FileContent(stateFile, DoesNotContain("- files"))
		t.FileSystem().FileContent(stateFile, Contains("- branches"))
	},
})
