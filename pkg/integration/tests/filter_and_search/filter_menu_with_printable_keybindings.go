package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterMenuWithPrintableKeybindings = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Driving a menu that filters as you type when every key configured for it is printable",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Keybinding.Universal.ConfirmMenu = config.Keybinding{"x"}
		cfg.GetUserConfig().Keybinding.Universal.Return = config.Keybinding{"q"}
		cfg.GetUserConfig().Keybinding.Universal.PrevItem = config.Keybinding{"k"}
		cfg.GetUserConfig().Keybinding.Universal.NextItem = config.Keybinding{"j"}
		cfg.GetUserConfig().Keybinding.Universal.PrevPage = config.Keybinding{"u"}
		cfg.GetUserConfig().Keybinding.Universal.NextPage = config.Keybinding{"d"}
		cfg.GetUserConfig().Keybinding.Universal.GotoTop = config.Keybinding{"g"}
		cfg.GetUserConfig().Keybinding.Universal.GotoBottom = config.Keybinding{"G"}
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		navigates := func(forward string, back string) {
			t.GlobalPress(config.Keybinding{forward})
			t.Views().Menu().SelectedLineIdxAtLeast(2)
			t.GlobalPress(config.Keybinding{back})
			t.Views().Menu().SelectedLineIdx(1)
		}

		t.Views().Files().IsFocused().Press(keys.Universal.OptionMenu)
		t.Views().Menu().IsFocused().SelectedLineIdx(1)

		// Until there is a filter, the configured keys drive the menu
		navigates("j", "k")
		navigates("d", "u")
		navigates("G", "g")

		// Once there is one, they are all filter text. It takes a key that isn't a
		// navigation key to get there.
		t.ExpectPopup().Menu().Filter("a")
		t.GlobalPress(config.Keybinding{"j"})
		t.GlobalPress(config.Keybinding{"k"})
		t.GlobalPress(config.Keybinding{"d"})
		t.GlobalPress(config.Keybinding{"u"})
		t.Views().MenuFilter().Content(Equals("ajkdu"))
		t.GlobalPress(config.Keybinding{"<c-u>"})

		// The menu is still navigable, because the physical keys drive it whatever
		// the configuration says
		navigates("<down>", "<up>")
		navigates("<pgdown>", "<pgup>")
		navigates("<end>", "<home>")

		// And so are confirming and cancelling. Escape gives up the filter first.
		t.ExpectPopup().Menu().Filter("Toggle whitespace")
		t.GlobalPress(config.Keybinding{"<esc>"})
		t.Views().MenuFilter().IsInvisible()
		t.Views().Menu().IsFocused().SelectedLine(Contains("Toggle whitespace"))

		t.ExpectPopup().Menu().Filter("Toggle whitespace")
		t.Views().Menu().SelectedLine(Contains("Toggle whitespace"))
		t.GlobalPress(config.Keybinding{"<enter>"})
		t.Views().Files().IsFocused()
	},
})
