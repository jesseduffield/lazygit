package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterMenuAsYouType = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filtering a menu by typing into the filter row that appears as you type",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo:    func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().Press(keys.Universal.OptionMenu)

		// The menu offers the filter, but stays as it is until we take it up on it
		t.Views().Menu().
			IsFocused().
			Subtitle(Equals("(Type to filter)")).
			Footer(Contains(" of "))
		t.Views().MenuFilter().IsInvisible()
		t.Views().MenuFilterFrame().IsInvisible()
		t.CursorIsHidden()

		t.ExpectPopup().Menu().Filter("whitespace")

		t.Views().Menu().
			Lines(
				Contains("─── Global"),
				Contains("Toggle whitespace").IsSelected(),
			).
			// the row covers the border the footer was on, so it moves there
			Subtitle(Equals("")).
			Footer(Equals(""))
		t.Views().MenuFilterFrame().
			IsVisible().
			Content(Equals("Filter ('@' for keybindings): ")).
			Footer(Equals("1 of 1")).
			SharesTopBorderWithBottomOf(t.Views().Menu())
		t.Views().MenuFilter().IsVisible().Content(Equals("whitespace"))
		t.Views().Tooltip().
			IsVisible().
			Content(Contains("Toggle whether or not whitespace changes are shown")).
			IsImmediatelyBelow(t.Views().MenuFilterFrame())
		t.CursorIsVisible()

		// Emptying the filter shows all the items again, and keeps the row
		t.GlobalPress(config.Keybinding{"<c-u>"})
		t.Views().MenuFilter().IsVisible().Content(Equals(""))
		t.Views().Menu().LineCount(GreaterThan(2))
		t.CursorIsVisible()

		// Moving the text cursor within the filter leaves the menu's selection alone
		t.ExpectPopup().Menu().Filter("co")
		t.Views().Menu().LineCount(GreaterThan(2))
		t.GlobalPress(config.Keybinding{"<down>"})
		t.Views().Menu().SelectedLineIdxAtLeast(2)
		t.GlobalPress(config.Keybinding{"<left>"})
		t.Views().Menu().SelectedLineIdxAtLeast(2)
		t.GlobalPress(config.Keybinding{"<right>"})

		// Clicking an item selects it and leaves the filter where it is
		t.Views().Menu().Click(0, 1).SelectedLineIdx(1)
		t.GlobalPress(config.Keybinding{"m"})
		t.Views().MenuFilter().Content(Equals("com"))

		t.GlobalPress(config.Keybinding{"<c-u>"})

		// Escape gives up the filter, keeping the item that was selected
		t.ExpectPopup().Menu().Filter("whitespace")
		t.Views().Menu().SelectedLine(Contains("Toggle whitespace"))
		t.GlobalPress(keys.Universal.Return)
		t.Views().Menu().
			IsFocused().
			SelectedLine(Contains("Toggle whitespace")).
			Subtitle(Equals("(Type to filter)")).
			Footer(Contains(" of "))
		t.Views().MenuFilter().IsInvisible()
		t.Views().MenuFilterFrame().IsInvisible()
		t.CursorIsHidden()

		// The next escape closes the menu
		t.GlobalPress(keys.Universal.Return)
		t.Views().Files().IsFocused()

		// A menu opened afterwards starts with no filter
		t.Views().Files().Press(keys.Universal.OptionMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			LineCount(GreaterThan(2)).
			Cancel()
	},
})
