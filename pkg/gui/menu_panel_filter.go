package gui

import (
	"github.com/jesseduffield/gocui"
	"regexp"
	"strings"
)

type MenuPanelFilter struct {
	needle     string
	menuView   *gocui.View
	menuItems  []*menuItem
	updateMenu func(string, []*menuItem)
}

func (m *MenuPanelFilter) applyFilterIfMatching(needle string) (newNeedle string) {
	filteredItems, _ := filterListItems(m.menuItems, needle)

	if len(filteredItems) > 0 {
		m.needle = needle
		m.updateMenu(needle, filteredItems)
	}

	newNeedle = m.needle
	return
}

func (m *MenuPanelFilter) HandleSearchKeystroke(key string) {
	isPrintableCharacter, _ := regexp.MatchString("[a-zA-Z]", key)

	if isPrintableCharacter {
		needle := m.needle + key
		m.applyFilterIfMatching(needle)
	}
}

func (m *MenuPanelFilter) HandleSearchBackspace() {
	if len(m.needle) > 0 {
		needle := m.needle[:len(m.needle)-1]
		m.applyFilterIfMatching(needle)
	}
}

func (m *MenuPanelFilter) HandleResetSearch() {
	m.needle = ""
	m.applyFilterIfMatching("")
}

func filterListItems(menuItems []*menuItem, filter string) ([]*menuItem, error) {
	var filteredItems []*menuItem
	for _, menuItem := range menuItems {
		for _, displayString := range menuItem.displayStrings {
			if strings.Contains(strings.ToUpper(displayString), strings.ToUpper(filter)) || len(filter) == 0 {
				filteredItems = append(filteredItems, menuItem)
				break
			}
		}
	}

	return filteredItems, nil
}
