package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushStackedBranchesShortcutConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Keep primary navigation bindings when they conflict with stacked push shortcuts",
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Keybinding.Universal.NextItem = []string{"a"}
	},
	SetupRepo: createRebasedStackOfBranches,
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().Press(keys.Universal.Push)

		t.ExpectPopup().Menu().
			Title(Equals("Push")).
			ContainsLines(
				Contains("Push all these branches in addition to the current one").DoesNotContain("a Push"),
				Contains("o Push only 'branch3'"),
			)

		t.Views().Menu().
			SelectedLine(Contains("Push all these branches in addition to the current one")).
			Press(keys.Universal.NextItem).
			SelectedLine(Contains("Push only 'branch3'")).
			PressEscape()

		t.Views().Files().IsFocused()
	},
})
