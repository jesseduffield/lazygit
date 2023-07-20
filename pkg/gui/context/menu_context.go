package context

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type MenuContext struct {
	c *ContextCommon

	*MenuViewModel
	*ListContextTrait
}

var _ types.IListContext = (*MenuContext)(nil)

func NewMenuContext(
	c *ContextCommon,
) *MenuContext {
	viewModel := NewMenuViewModel(c)

	return &MenuContext{
		c:             c,
		MenuViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:                  c.Views().Menu,
				WindowName:            "menu",
				Key:                   "menu",
				Kind:                  types.TEMPORARY_POPUP,
				Focusable:             true,
				HasUncontrolledBounds: true,
			})),
			getDisplayStrings:   viewModel.GetDisplayStrings,
			list:                viewModel,
			c:                   c,
			getColumnAlignments: func() []utils.Alignment { return viewModel.columnAlignment },
		},
	}
}

// TODO: remove this thing.
func (self *MenuContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.Label
}

type MenuViewModel struct {
	c               *ContextCommon
	menuItems       []*types.MenuItem
	columnAlignment []utils.Alignment
	*FilteredListViewModel[*types.MenuItem]
}

func NewMenuViewModel(c *ContextCommon) *MenuViewModel {
	self := &MenuViewModel{
		menuItems: nil,
		c:         c,
	}

	self.FilteredListViewModel = NewFilteredListViewModel(
		func() []*types.MenuItem { return self.menuItems },
		func(item *types.MenuItem) []string { return item.LabelColumns },
	)

	return self
}

func (self *MenuViewModel) SetMenuItems(items []*types.MenuItem, columnAlignment []utils.Alignment) {
	self.menuItems = items
	self.columnAlignment = columnAlignment
}

// TODO: move into presentation package
func (self *MenuViewModel) GetDisplayStrings(_startIdx int, _length int) [][]string {
	menuItems := self.FilteredListViewModel.GetItems()
	showKeys := slices.Some(menuItems, func(item *types.MenuItem) bool {
		return item.Key != nil
	})

	return slices.Map(menuItems, func(item *types.MenuItem) []string {
		displayStrings := item.LabelColumns

		if !showKeys {
			return displayStrings
		}

		// These keys are used for general navigation so we'll strike them out to
		// avoid confusion
		reservedKeys := []string{
			self.c.UserConfig.Keybinding.Universal.Confirm,
			self.c.UserConfig.Keybinding.Universal.Select,
			self.c.UserConfig.Keybinding.Universal.Return,
			self.c.UserConfig.Keybinding.Universal.StartSearch,
		}
		keyLabel := keybindings.LabelFromKey(item.Key)
		keyStyle := style.FgCyan
		if lo.Contains(reservedKeys, keyLabel) {
			keyStyle = style.FgDefault.SetStrikethrough()
		}

		displayStrings = slices.Prepend(displayStrings, keyStyle.Sprint(keyLabel))
		return displayStrings
	})
}

func (self *MenuContext) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	basicBindings := self.ListContextTrait.GetKeybindings(opts)
	menuItemsWithKeys := slices.Filter(self.menuItems, func(item *types.MenuItem) bool {
		return item.Key != nil
	})

	menuItemBindings := slices.Map(menuItemsWithKeys, func(item *types.MenuItem) *types.Binding {
		return &types.Binding{
			Key:     item.Key,
			Handler: func() error { return self.OnMenuPress(item) },
		}
	})

	// appending because that means the menu item bindings have lower precedence.
	// So if a basic binding is to escape from the menu, we want that to still be
	// what happens when you press escape. This matters when we're showing the menu
	// for all keybindings of say the files context.
	return append(basicBindings, menuItemBindings...)
}

func (self *MenuContext) OnMenuPress(selectedItem *types.MenuItem) error {
	if err := self.c.PopContext(); err != nil {
		return err
	}

	if selectedItem == nil {
		return nil
	}

	if err := selectedItem.OnPress(); err != nil {
		return err
	}

	return nil
}
