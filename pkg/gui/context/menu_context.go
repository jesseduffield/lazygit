package context

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MenuContext struct {
	*MenuViewModel
	*ListContextTrait
}

var _ types.IListContext = (*MenuContext)(nil)

func NewMenuContext(
	view *gocui.View,

	c *types.HelperCommon,
	getOptionsMap func() map[string]string,
	renderToDescriptionView func(string),
) *MenuContext {
	viewModel := NewMenuViewModel()

	onFocus := func(types.OnFocusOpts) error {
		selectedMenuItem := viewModel.GetSelected()
		renderToDescriptionView(selectedMenuItem.Tooltip)
		return nil
	}

	return &MenuContext{
		MenuViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:                  view,
				WindowName:            "menu",
				Key:                   "menu",
				Kind:                  types.TEMPORARY_POPUP,
				OnGetOptionsMap:       getOptionsMap,
				Focusable:             true,
				HasUncontrolledBounds: true,
			}), ContextCallbackOpts{
				OnFocus: onFocus,
			}),
			getDisplayStrings: viewModel.GetDisplayStrings,
			list:              viewModel,
			c:                 c,
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
	menuItems []*types.MenuItem
	*BasicViewModel[*types.MenuItem]
}

func NewMenuViewModel() *MenuViewModel {
	self := &MenuViewModel{
		menuItems: nil,
	}

	self.BasicViewModel = NewBasicViewModel(func() []*types.MenuItem { return self.menuItems })

	return self
}

func (self *MenuViewModel) SetMenuItems(items []*types.MenuItem) {
	self.menuItems = items
}

// TODO: move into presentation package
func (self *MenuViewModel) GetDisplayStrings(_startIdx int, _length int) [][]string {
	showKeys := slices.Some(self.menuItems, func(item *types.MenuItem) bool {
		return item.Key != nil
	})

	return slices.Map(self.menuItems, func(item *types.MenuItem) []string {
		displayStrings := item.LabelColumns
		if showKeys {
			displayStrings = slices.Prepend(displayStrings, style.FgCyan.Sprint(keybindings.LabelFromKey(item.Key)))
		}
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

	if err := selectedItem.OnPress(); err != nil {
		return err
	}

	return nil
}
