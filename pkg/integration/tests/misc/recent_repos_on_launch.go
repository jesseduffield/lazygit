package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Couldn't find an easy way to actually reproduce the situation of opening outside a repo,
// so I'm introducing a hacky env var to force lazygit to show the recent repos menu upon opening.

var RecentReposOnLaunch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When opening to a menu, focus is correctly given to the menu",
	ExtraCmdArgs: []string{},
	ExtraEnvVars: map[string]string{
		"SHOW_RECENT_REPOS": "true",
	},
	Skip:        false,
	SetupConfig: func(config *config.AppConfig) {},
	SetupRepo:   func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectPopup().Menu().
			Title(Equals("Recent repositories")).
			Select(Contains("Cancel")).
			Confirm()

		t.Views().Files().IsFocused()
	},
})
