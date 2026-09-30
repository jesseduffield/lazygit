package helpers

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// ThemeHelper lets the user switch between the theme files in the themes
// folder of the config dir.
type ThemeHelper struct {
	c               *HelperCommon
	onThemeSelected func()
}

// NewThemeHelper takes onThemeSelected, which is called after the theme
// settings of the user config have changed, to apply them and to render again
// the views whose content was styled with the previous colors.
func NewThemeHelper(c *HelperCommon, onThemeSelected func()) *ThemeHelper {
	return &ThemeHelper{
		c:               c,
		onThemeSelected: onThemeSelected,
	}
}

func (self *ThemeHelper) OpenThemeMenu() error {
	appConfig := self.c.GetConfig()
	themes, err := appConfig.ListThemes()
	if err != nil {
		return err
	}
	selectedTheme := appConfig.GetSelectedTheme()
	appliedTheme := appConfig.GetAppliedTheme()

	// The error of the last failed attempt to load the selected theme goes
	// into the tooltip of the selected item. That is "(none)" when the file
	// that says which theme is selected couldn't be read.
	loadErrorText := ""
	if loadError := appConfig.GetThemeLoadError(); loadError != nil {
		loadErrorText = loadError.Error()
	}
	themeItem := func(themeName string, label string) *types.MenuItem {
		item := &types.MenuItem{
			Label:   label,
			OnPress: func() error { return self.selectTheme(themeName) },
			Widget:  types.MakeMenuRadioButton(themeName == selectedTheme),
		}
		if themeName == selectedTheme {
			item.Tooltip = loadErrorText
		}
		return item
	}

	listedThemeItems := lo.Map(themes, func(themeName string, _ int) *types.MenuItem {
		label := themeName
		// The selected theme's file couldn't be loaded, or didn't exist yet,
		// when the config was last loaded; choosing the item loads it again
		if themeName == selectedTheme && themeName != appliedTheme {
			label = utils.ResolvePlaceholderString(self.c.Tr.ThemeNotLoadedLabel, map[string]string{
				"name": themeName,
			})
		}
		return themeItem(themeName, label)
	})
	menuItems := append([]*types.MenuItem{themeItem("", self.c.Tr.NoTheme)}, listedThemeItems...)

	// A selected theme that isn't listed, usually because its file is
	// missing, isn't applied, but it stays selected so that it comes back
	// with its file. Show it, so that the menu doesn't pretend that no theme
	// is selected.
	if selectedTheme != "" && !lo.Contains(themes, selectedTheme) {
		missingThemeLabel := utils.ResolvePlaceholderString(self.c.Tr.MissingTheme, map[string]string{
			"name": selectedTheme,
		})
		missingThemeItem := themeItem(selectedTheme, missingThemeLabel)
		missingThemeItem.DisabledReason = &types.DisabledReason{Text: self.themeNotFoundMessage(selectedTheme)}
		menuItems = append(menuItems, missingThemeItem)
	}

	prompt := ""
	if len(themes) == 0 {
		prompt = utils.ResolvePlaceholderString(self.c.Tr.NoThemesFound, map[string]string{
			"dir": appConfig.GetThemesDir(),
		})
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title:           self.c.Tr.SelectTheme,
		Prompt:          prompt,
		Items:           menuItems,
		FilterAsYouType: true,
	})
}

func (self *ThemeHelper) selectTheme(name string) error {
	if err := self.c.GetConfig().SelectTheme(name); err != nil {
		if errors.Is(err, config.ErrThemeNotFound) {
			return errors.New(self.themeNotFoundMessage(name))
		}
		return err
	}

	self.onThemeSelected()

	// The contexts of the staging, patch building and merge conflicts views
	// render them when they are focused, and applying the theme renders the
	// main view again only for a side panel. Closing the theme menu has
	// already focused the current context again, before the theme was
	// applied, and the renders started then finish in the background, so they
	// would bring back the previous theme's colors. Activating it once more
	// renders them again after those, with the new colors; a side panel's
	// main view has been rendered with them already, so it is left alone. It
	// also renders the search status again, whose frame color applying the
	// theme has reset.
	self.c.Context().Activate(self.c.Context().Current(), types.OnFocusOpts{SkipMainViewUpdate: true})

	displayName := name
	if name == "" {
		displayName = self.c.Tr.NoTheme
	}
	self.c.Toast(utils.ResolvePlaceholderString(self.c.Tr.SelectedTheme, map[string]string{
		"name": displayName,
	}))
	return nil
}

func (self *ThemeHelper) themeNotFoundMessage(name string) string {
	return utils.ResolvePlaceholderString(self.c.Tr.ThemeNotFound, map[string]string{
		"name": name,
		"dir":  self.c.GetConfig().GetThemesDir(),
	})
}
