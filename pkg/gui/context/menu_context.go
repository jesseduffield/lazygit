package context

import (
	"errors"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
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
			ListRenderer: ListRenderer{
				list:                viewModel,
				getDisplayStrings:   viewModel.GetDisplayStrings,
				getColumnAlignments: func() []utils.Alignment { return viewModel.columnAlignment },
				getNonModelItems:    viewModel.GetNonModelItems,
			},
			c: c,
		},
	}
}

type MenuViewModel struct {
	c                         *ContextCommon
	menuItems                 []*types.MenuItem
	prompt                    string
	promptLines               []string
	columnAlignment           []utils.Alignment
	allowFilteringKeybindings bool
	keybindingsTakePrecedence bool
	*FilteredListViewModel[*types.MenuItem]
}

func NewMenuViewModel(c *ContextCommon) *MenuViewModel {
	self := &MenuViewModel{
		menuItems: nil,
		c:         c,
	}

	filterKeybindings := false

	self.FilteredListViewModel = NewFilteredListViewModel(
		func() []*types.MenuItem { return self.menuItems },
		func(item *types.MenuItem) []string {
			if filterKeybindings {
				return []string{keybindings.LabelFromKey(item.Key)}
			}

			return item.LabelColumns
		},
	)

	self.FilteredListViewModel.SetPreprocessFilterFunc(func(filter string) string {
		if self.allowFilteringKeybindings && strings.HasPrefix(filter, "@") {
			filterKeybindings = true
			return filter[1:]
		}

		filterKeybindings = false
		return filter
	})

	return self
}

func (self *MenuViewModel) SetMenuItems(items []*types.MenuItem, columnAlignment []utils.Alignment) {
	self.menuItems = items
	self.columnAlignment = columnAlignment
}

func (self *MenuViewModel) GetPrompt() string {
	return self.prompt
}

func (self *MenuViewModel) SetPrompt(prompt string) {
	self.prompt = prompt
	self.promptLines = nil
}

func (self *MenuViewModel) GetPromptLines() []string {
	return self.promptLines
}

func (self *MenuViewModel) SetPromptLines(promptLines []string) {
	self.promptLines = promptLines
}

func (self *MenuViewModel) SetAllowFilteringKeybindings(allow bool) {
	self.allowFilteringKeybindings = allow
}

func (self *MenuViewModel) SetKeybindingsTakePrecedence(value bool) {
	self.keybindingsTakePrecedence = value
}

// TODO: move into presentation package
func (self *MenuViewModel) GetDisplayStrings(_ int, _ int) [][]string {
	menuItems := self.FilteredListViewModel.GetItems()

	return lo.Map(menuItems, func(item *types.MenuItem, _ int) []string {
		displayStrings := item.LabelColumns
		if item.DisabledReason != nil {
			displayStrings[0] = style.FgDefault.SetStrikethrough().Sprint(displayStrings[0])
		}

		keyLabel := ""
		if item.Key != nil {
			keyLabel = style.FgCyan.Sprint(keybindings.LabelFromKey(item.Key))
		}

		checkMark := ""
		switch item.Widget {
		case types.MenuWidgetNone:
			// do nothing
		case types.MenuWidgetRadioButtonSelected:
			checkMark = "(•)"
		case types.MenuWidgetRadioButtonUnselected:
			checkMark = "( )"
		case types.MenuWidgetCheckboxSelected:
			checkMark = "[✓]"
		case types.MenuWidgetCheckboxUnselected:
			checkMark = "[ ]"
		}

		displayStrings = utils.Prepend(displayStrings, keyLabel, checkMark)
		return displayStrings
	})
}

func (self *MenuViewModel) GetNonModelItems() []*NonModelItem {
	result := []*NonModelItem{}
	result = append(result, lo.Map(self.promptLines, func(line string, _ int) *NonModelItem {
		return &NonModelItem{
			Index:   0,
			Column:  0,
			Content: line,
		}
	})...)

	// Don't display section headers when we are filtering, and the filter mode
	// is fuzzy. The reason is that filtering changes the order of the items
	// (they are sorted by best match), so all the sections would be messed up.
	if self.FilteredListViewModel.IsFiltering() && self.c.UserConfig().Gui.UseFuzzySearch() {
		return result
	}

	menuItems := self.FilteredListViewModel.GetItems()
	var prevSection *types.MenuSection
	for i, menuItem := range menuItems {
		if menuItem.Section != nil && menuItem.Section != prevSection {
			if prevSection != nil {
				result = append(result, &NonModelItem{
					Index:   i,
					Column:  1,
					Content: "",
				})
			}

			result = append(result, &NonModelItem{
				Index:   i,
				Column:  1,
				Content: style.FgGreen.SetBold().Sprintf("--- %s ---", menuItem.Section.Title),
			})
			prevSection = menuItem.Section
		}
	}

	return result
}

func (self *MenuContext) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	basicBindings := self.ListContextTrait.GetKeybindings(opts)
	menuItemsWithKeys := lo.Filter(self.menuItems, func(item *types.MenuItem, _ int) bool {
		return item.Key != nil
	})

	menuItemBindings := lo.Map(menuItemsWithKeys, func(item *types.MenuItem, _ int) *types.Binding {
		return &types.Binding{
			Key:     item.Key,
			Handler: func() error { return self.OnMenuPress(item) },
		}
	})

	if self.keybindingsTakePrecedence {
		// This is used for all normal menus except the keybindings menu. In this case we want the
		// bindings of the menu items to have higher precedence than the builtin bindings; this
		// allows assigning a keybinding to a menu item that overrides a non-essential binding such
		// as 'j', 'k', 'H', 'L', etc. This is safe to do because the essential bindings such as
		// confirm and return have already been removed from the menu items in this case.
		return append(menuItemBindings, basicBindings...)
	}

	// For the keybindings menu we didn't remove the essential bindings from the menu items, because
	// it is important to see all bindings (as a cheat sheet for what the keys are when the menu is
	// not open). Therefore we want the essential bindings to have higher precedence than the menu
	// item bindings.
	return append(basicBindings, menuItemBindings...)
}

func (self *MenuContext) OnMenuPress(selectedItem *types.MenuItem) error {
	if selectedItem != nil && selectedItem.DisabledReason != nil {
		if selectedItem.DisabledReason.ShowErrorInPanel {
			return errors.New(selectedItem.DisabledReason.Text)
		}

		self.c.ErrorToast(self.c.Tr.DisabledMenuItemPrefix + selectedItem.DisabledReason.Text)
		return nil
	}

	self.c.Context().Pop()

	if selectedItem == nil {
		return nil
	}

	if err := selectedItem.OnPress(); err != nil {
		return err
	}

	return nil
}

// There is currently no need to use range-select in a menu so we're disabling it.
func (self *MenuContext) RangeSelectEnabled() bool {
	return false
}

func (self *MenuContext) FilterPrefix(tr *i18n.TranslationSet) string {
	if self.allowFilteringKeybindings {
		return tr.FilterPrefixMenu
	}

	return self.FilteredListViewModel.FilterPrefix(tr)
}
