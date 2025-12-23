package helpers

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type GitConfigHelper struct {
	c                *HelperCommon
	displayScope     gitConfigDisplayScope
	expandedSections map[string]bool
}

func NewGitConfigHelper(c *HelperCommon) *GitConfigHelper {
	expandedSections := map[string]bool{}
	for _, section := range c.GetAppState().GitConfigExpandedSections {
		expandedSections[section] = true
	}
	return &GitConfigHelper{
		c:                c,
		displayScope:     gitConfigDisplayScopeGlobal,
		expandedSections: expandedSections,
	}
}

func (self *GitConfigHelper) OpenMenu() error {
	return self.openMenuWithSelection(self.currentMenuSelectionID())
}

func (self *GitConfigHelper) openMenuWithSelection(selectionID string) error {
	localConfig := self.c.Git().Config.ListLocalConfig()
	globalConfig := self.c.Git().Config.ListGlobalConfig()
	systemConfig := self.c.Git().Config.ListSystemConfig()

	displayValues := self.displayScopeValues(localConfig, globalConfig, systemConfig)
	keys := keysForScope(displayValues, defaultGitConfigKeysForScope(self.displayScope))
	configSection := &types.MenuSection{
		Title: utils.ResolvePlaceholderString(
			self.c.Tr.GitConfigConfigSectionWithScope,
			map[string]string{"scope": self.displayScopeLabel()},
		),
		Column: 0,
	}
	menuItems := []*types.MenuItem{}

	menuItems = append(menuItems, self.sectionMenuItems(configSection, keys, displayValues)...)

	self.c.Contexts().Menu.SetExtraKeybindings(self.scopeKeybindings())

	if err := self.c.Menu(types.CreateMenuOptions{
		Title:           self.c.Tr.GitConfigTitle,
		Items:           menuItems,
		ColumnAlignment: []utils.Alignment{utils.AlignLeft, utils.AlignLeft},
	}); err != nil {
		return err
	}
	self.restoreSelection(menuItems, selectionID)
	return nil
}

func (self *GitConfigHelper) openGitConfigEntryMenu(key string) error {
	localValue := self.c.Git().Config.GetLocalConfigValue(key)
	globalValue := self.c.Git().Config.GetGlobalConfigValue(key)
	notSetReason := &types.DisabledReason{Text: self.c.Tr.GitConfigNotSet}

	menuItems := []*types.MenuItem{
		{
			LabelColumns: []string{
				self.c.Tr.GitConfigSetLocal,
				self.formatGitConfigValue(localValue, style.FgGreen),
			},
			OnPress: func() error {
				return self.promptSetGitConfigValue(key, gitConfigScopeLocal, localValue)
			},
		},
		{
			LabelColumns: []string{
				self.c.Tr.GitConfigSetGlobal,
				self.formatGitConfigValue(globalValue, style.FgYellow),
			},
			OnPress: func() error {
				return self.promptSetGitConfigValue(key, gitConfigScopeGlobal, globalValue)
			},
		},
		{
			Label:          self.c.Tr.GitConfigUnsetLocal,
			OnPress:        func() error { return self.unsetGitConfigValue(key, gitConfigScopeLocal) },
			DisabledReason: lo.Ternary(localValue == "", notSetReason, nil),
		},
		{
			Label:          self.c.Tr.GitConfigUnsetGlobal,
			OnPress:        func() error { return self.unsetGitConfigValue(key, gitConfigScopeGlobal) },
			DisabledReason: lo.Ternary(globalValue == "", notSetReason, nil),
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title:           key,
		Items:           menuItems,
		ColumnAlignment: []utils.Alignment{utils.AlignLeft, utils.AlignLeft},
	})
}

type gitConfigScope struct {
	name string
}

var (
	gitConfigScopeLocal  = gitConfigScope{name: "local"}
	gitConfigScopeGlobal = gitConfigScope{name: "global"}
)

func (self *GitConfigHelper) promptSetGitConfigValue(key string, scope gitConfigScope, currentValue string) error {
	title := utils.ResolvePlaceholderString(
		self.c.Tr.GitConfigSetPromptTitle,
		map[string]string{
			"key":   key,
			"scope": self.gitConfigScopeLabel(scope),
		},
	)
	self.c.Prompt(types.PromptOpts{
		Title:              title,
		InitialContent:     currentValue,
		AllowEmptyInput:    false,
		PreserveWhitespace: true,
		HandleConfirm: func(value string) error {
			if err := self.setGitConfigValue(key, scope, value); err != nil {
				return err
			}
			self.c.Toast(utils.ResolvePlaceholderString(self.c.Tr.GitConfigUpdatedToast, map[string]string{
				"key":   key,
				"scope": self.gitConfigScopeLabel(scope),
			}))
			return self.OpenMenu()
		},
	})
	return nil
}

func (self *GitConfigHelper) setGitConfigValue(key string, scope gitConfigScope, value string) error {
	var err error
	switch scope.name {
	case "local":
		err = self.c.Git().Config.SetLocalConfigValue(key, value)
	case "global":
		err = self.c.Git().Config.SetGlobalConfigValue(key, value)
	default:
		return errors.New("unknown git config scope")
	}

	if err != nil {
		return err
	}

	self.c.Git().Config.DropConfigCache()
	return nil
}

func (self *GitConfigHelper) unsetGitConfigValue(key string, scope gitConfigScope) error {
	var err error
	switch scope.name {
	case "local":
		err = self.c.Git().Config.UnsetLocalConfigValue(key)
	case "global":
		err = self.c.Git().Config.UnsetGlobalConfigValue(key)
	default:
		return errors.New("unknown git config scope")
	}

	if err != nil {
		return err
	}

	self.c.Git().Config.DropConfigCache()
	self.c.Toast(utils.ResolvePlaceholderString(self.c.Tr.GitConfigUpdatedToast, map[string]string{
		"key":   key,
		"scope": self.gitConfigScopeLabel(scope),
	}))
	return self.OpenMenu()
}

func (self *GitConfigHelper) formatGitConfigValue(value string, color style.TextStyle) string {
	if value == "" {
		return style.FgBlackLighter.Sprint(self.c.Tr.GitConfigNotSet)
	}
	return color.Sprint(value)
}

func (self *GitConfigHelper) gitConfigScopeLabel(scope gitConfigScope) string {
	switch scope.name {
	case "local":
		return self.c.Tr.GitConfigScopeLocal
	case "global":
		return self.c.Tr.GitConfigScopeGlobal
	default:
		return scope.name
	}
}

type gitConfigDisplayScope struct {
	name string
}

var (
	gitConfigDisplayScopeLocal  = gitConfigDisplayScope{name: "local"}
	gitConfigDisplayScopeGlobal = gitConfigDisplayScope{name: "global"}
	gitConfigDisplayScopeSystem = gitConfigDisplayScope{name: "system"}
)

func (self gitConfigDisplayScope) label(tr *i18n.TranslationSet) string {
	switch self.name {
	case "local":
		return tr.GitConfigLocalColumn
	case "global":
		return tr.GitConfigGlobalColumn
	case "system":
		return tr.GitConfigSystemColumn
	default:
		return self.name
	}
}

func (self gitConfigDisplayScope) keybinding() types.Key {
	switch self.name {
	case "local":
		return 'l'
	case "global":
		return 'g'
	case "system":
		return 's'
	default:
		return nil
	}
}

func (self *GitConfigHelper) displayScopeValues(localConfig map[string]string, globalConfig map[string]string, systemConfig map[string]string) map[string]string {
	switch self.displayScope.name {
	case "local":
		return localConfig
	case "global":
		return globalConfig
	case "system":
		return systemConfig
	default:
		return map[string]string{}
	}
}

func (self *GitConfigHelper) sectionMenuItems(sectionHeader *types.MenuSection, keys []string, values map[string]string) []*types.MenuItem {
	sections := map[string][]string{}
	leafKeys := []string{}
	for _, key := range keys {
		parts := strings.SplitN(key, ".", 2)
		if len(parts) == 1 {
			leafKeys = append(leafKeys, key)
			continue
		}
		sections[parts[0]] = append(sections[parts[0]], parts[1])
	}

	sectionNames := make([]string, 0, len(sections))
	for sectionName := range sections {
		sectionNames = append(sectionNames, sectionName)
	}
	sort.Strings(sectionNames)
	sort.Strings(leafKeys)

	menuItems := make([]*types.MenuItem, 0, len(keys))
	for _, section := range sectionNames {
		isExpanded := self.isSectionExpanded(section)
		prefix := "+ "
		if isExpanded {
			prefix = "- "
		}
		menuItems = append(menuItems, &types.MenuItem{
			Label:        "section:" + section,
			LabelColumns: []string{prefix + section},
			Section:      sectionHeader,
			OnPress: func() error {
				self.toggleSection(section)
				return self.openMenuWithSelection("section:" + section)
			},
		})

		if !isExpanded {
			continue
		}

		children := sections[section]
		sort.Strings(children)
		for i, child := range children {
			isLast := i == len(children)-1
			branch := "  |- "
			if isLast {
				branch = "  `- "
			}
			fullKey := section + "." + child
			menuItems = append(menuItems, &types.MenuItem{
				Label: "key:" + fullKey,
				LabelColumns: []string{
					branch + child,
					self.formatGitConfigValue(values[fullKey], self.displayScopeColor()),
				},
				Section: sectionHeader,
				OnPress: func() error {
					return self.openGitConfigEntryMenu(fullKey)
				},
			})
		}
	}

	for _, key := range leafKeys {
		menuItems = append(menuItems, &types.MenuItem{
			Label: "key:" + key,
			LabelColumns: []string{
				key,
				self.formatGitConfigValue(values[key], self.displayScopeColor()),
			},
			Section: sectionHeader,
			OnPress: func() error {
				return self.openGitConfigEntryMenu(key)
			},
		})
	}

	return menuItems
}

func (self *GitConfigHelper) displayScopeColor() style.TextStyle {
	switch self.displayScope.name {
	case "local":
		return style.FgGreen
	case "global":
		return style.FgYellow
	case "system":
		return style.FgBlue
	default:
		return style.FgDefault
	}
}

func (self *GitConfigHelper) isSectionExpanded(section string) bool {
	expanded, ok := self.expandedSections[section]
	return ok && expanded
}

func (self *GitConfigHelper) toggleSection(section string) {
	self.expandedSections[section] = !self.isSectionExpanded(section)
	self.persistExpandedSections()
}

func (self *GitConfigHelper) persistExpandedSections() {
	expanded := []string{}
	for section, isExpanded := range self.expandedSections {
		if isExpanded {
			expanded = append(expanded, section)
		}
	}
	sort.Strings(expanded)
	self.c.GetAppState().GitConfigExpandedSections = expanded
	self.c.SaveAppStateAndLogError()
}

func (self *GitConfigHelper) currentMenuSelectionID() string {
	if self.c.Contexts().Menu == nil {
		return ""
	}
	return self.c.Contexts().Menu.GetSelectedItemId()
}

func (self *GitConfigHelper) restoreSelection(items []*types.MenuItem, selectionID string) {
	if selectionID == "" {
		return
	}
	for i, item := range items {
		if item != nil && item.ID() == selectionID {
			self.c.Contexts().Menu.GetList().SetSelection(i)
			self.c.Contexts().Menu.FocusLine(true)
			return
		}
	}
}

func (self *GitConfigHelper) scopeKeybindings() []*types.Binding {
	return []*types.Binding{
		{
			Key:     gitConfigDisplayScopeLocal.keybinding(),
			Handler: func() error { return self.setDisplayScope(gitConfigDisplayScopeLocal) },
		},
		{
			Key:     gitConfigDisplayScopeGlobal.keybinding(),
			Handler: func() error { return self.setDisplayScope(gitConfigDisplayScopeGlobal) },
		},
		{
			Key:     gitConfigDisplayScopeSystem.keybinding(),
			Handler: func() error { return self.setDisplayScope(gitConfigDisplayScopeSystem) },
		},
		{
			Key:     gocui.KeyArrowLeft,
			Handler: self.selectPreviousScope,
		},
		{
			Key:     gocui.KeyArrowRight,
			Handler: self.selectNextScope,
		},
	}
}

func (self *GitConfigHelper) setDisplayScope(scope gitConfigDisplayScope) error {
	if self.displayScope.name == scope.name {
		return nil
	}
	self.displayScope = scope
	return self.openMenuWithSelection(self.currentMenuSelectionID())
}

func (self *GitConfigHelper) selectPreviousScope() error {
	return self.selectScopeOffset(-1)
}

func (self *GitConfigHelper) selectNextScope() error {
	return self.selectScopeOffset(1)
}

func (self *GitConfigHelper) selectScopeOffset(offset int) error {
	scopes := []gitConfigDisplayScope{
		gitConfigDisplayScopeLocal,
		gitConfigDisplayScopeGlobal,
		gitConfigDisplayScopeSystem,
	}
	currentIdx := 0
	for i, scope := range scopes {
		if scope.name == self.displayScope.name {
			currentIdx = i
			break
		}
	}
	nextIdx := (currentIdx + offset) % len(scopes)
	if nextIdx < 0 {
		nextIdx += len(scopes)
	}
	return self.setDisplayScope(scopes[nextIdx])
}

func (self *GitConfigHelper) displayScopeLabel() string {
	return fmt.Sprintf("< %s >", self.displayScope.label(self.c.Tr))
}

func keysForScope(values map[string]string, defaults []string) []string {
	keySet := make(map[string]struct{}, len(values)+len(defaults))
	for key := range values {
		keySet[key] = struct{}{}
	}
	for _, key := range defaults {
		keySet[key] = struct{}{}
	}
	result := make([]string, 0, len(keySet))
	for key := range keySet {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

type defaultGitConfigKey struct {
	key    string
	scopes map[string]bool
}

func defaultGitConfigKeysForScope(scope gitConfigDisplayScope) []string {
	keys := []defaultGitConfigKey{
		{key: "user.name", scopes: map[string]bool{"local": true, "global": true}},
		{key: "user.email", scopes: map[string]bool{"local": true, "global": true}},
		{key: "core.editor", scopes: map[string]bool{"local": true, "global": true, "system": true}},
		{key: "core.autocrlf", scopes: map[string]bool{"local": true, "global": true, "system": true}},
		{key: "core.ignorecase", scopes: map[string]bool{"local": true, "global": true, "system": true}},
		{key: "init.defaultBranch", scopes: map[string]bool{"global": true, "system": true}},
		{key: "pull.rebase", scopes: map[string]bool{"local": true, "global": true}},
		{key: "pull.ff", scopes: map[string]bool{"local": true, "global": true}},
		{key: "merge.ff", scopes: map[string]bool{"local": true, "global": true}},
		{key: "push.default", scopes: map[string]bool{"local": true, "global": true}},
		{key: "fetch.prune", scopes: map[string]bool{"local": true, "global": true}},
		{key: "commit.gpgSign", scopes: map[string]bool{"local": true, "global": true}},
		{key: "tag.gpgSign", scopes: map[string]bool{"local": true, "global": true}},
		{key: "gpg.program", scopes: map[string]bool{"global": true, "system": true}},
		{key: "rebase.autostash", scopes: map[string]bool{"local": true, "global": true}},
		{key: "rerere.enabled", scopes: map[string]bool{"local": true, "global": true}},
		{key: "diff.tool", scopes: map[string]bool{"local": true, "global": true}},
		{key: "merge.tool", scopes: map[string]bool{"local": true, "global": true}},
	}
	result := []string{}
	for _, entry := range keys {
		if entry.scopes[scope.name] {
			result = append(result, entry.key)
		}
	}
	sort.Strings(result)
	return result
}
