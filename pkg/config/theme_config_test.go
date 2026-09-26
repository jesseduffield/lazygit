package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThemeForBackground(t *testing.T) {
	gui := GuiConfig{
		Theme: ThemeConfig{
			ActiveBorderColor:   []string{"green"},
			InactiveBorderColor: []string{"default"},
			AuthorColors:        map[string]string{"Alice": "red", "Bob": "blue"},
			BranchColorPatterns: ColorPatterns{
				{Pattern: "^docs/", Color: "cyan"},
				{Pattern: "^feature/", Color: "green"},
			},
		},
		DarkTheme: ThemeConfig{
			ActiveBorderColor: []string{"yellow"},
		},
		LightTheme: ThemeConfig{
			InactiveBorderColor: []string{"#777777"},
			AuthorColors:        map[string]string{"Bob": "#000080"},
			BranchColorPatterns: ColorPatterns{
				{Pattern: "ISSUE", Color: "red"},
				{Pattern: "^docs/", Color: "#008080"},
			},
		},
	}

	dark := gui.ThemeForBackground(false, "")
	assert.Equal(t, []string{"yellow"}, dark.ActiveBorderColor)
	assert.Equal(t, []string{"default"}, dark.InactiveBorderColor)
	assert.Equal(t, map[string]string{"Alice": "red", "Bob": "blue"}, dark.AuthorColors)
	assert.Equal(t, gui.Theme.BranchColorPatterns, dark.BranchColorPatterns)

	light := gui.ThemeForBackground(true, "")
	assert.Equal(t, []string{"green"}, light.ActiveBorderColor)
	assert.Equal(t, []string{"#777777"}, light.InactiveBorderColor)
	assert.Equal(t, map[string]string{"Alice": "red", "Bob": "#000080"}, light.AuthorColors)
	assert.Equal(t, ColorPatterns{
		{Pattern: "ISSUE", Color: "red"},
		{Pattern: "^docs/", Color: "#008080"},
		{Pattern: "^feature/", Color: "green"},
	}, light.BranchColorPatterns)

	assert.Equal(t, map[string]string{"Alice": "red", "Bob": "blue"}, gui.Theme.AuthorColors,
		"merging must leave the themes it merges alone")
}

func TestEveryThemeFieldCanBeOverridden(t *testing.T) {
	var override ThemeConfig
	overrideValue := reflect.ValueOf(&override).Elem()
	for i := range overrideValue.NumField() {
		field := overrideValue.Field(i)
		switch field.Interface().(type) {
		case []string:
			field.Set(reflect.ValueOf([]string{"#123456"}))
		case map[string]string:
			field.Set(reflect.ValueOf(map[string]string{"key": "#123456"}))
		case ColorPatterns:
			field.Set(reflect.ValueOf(ColorPatterns{{Pattern: "key", Color: "#123456"}}))
		default:
			t.Fatalf("no test value for theme field %s", overrideValue.Type().Field(i).Name)
		}
	}

	gui := GetDefaultConfig().Gui
	gui.DarkTheme = override
	assert.Equal(t, override, gui.ThemeForBackground(false, ""))
}

func TestThemeForBackgroundFallsBackToTheDefaultsForTheBackground(t *testing.T) {
	gui := GetDefaultConfig().Gui
	assert.Equal(t, []string{"#4d4d4d"}, gui.ThemeForBackground(false, "").InactiveViewSelectedLineBgColor)
	assert.Equal(t, []string{"#d9d9d9"}, gui.ThemeForBackground(true, "").InactiveViewSelectedLineBgColor)

	gui.Theme.InactiveViewSelectedLineBgColor = []string{"bold"}
	assert.Equal(t, []string{"bold"}, gui.ThemeForBackground(false, "").InactiveViewSelectedLineBgColor)
	assert.Equal(t, []string{"bold"}, gui.ThemeForBackground(true, "").InactiveViewSelectedLineBgColor)

	gui.LightTheme.InactiveViewSelectedLineBgColor = []string{"white"}
	assert.Equal(t, []string{"bold"}, gui.ThemeForBackground(false, "").InactiveViewSelectedLineBgColor)
	assert.Equal(t, []string{"white"}, gui.ThemeForBackground(true, "").InactiveViewSelectedLineBgColor)
}

func TestBackgroundDefaultsAreDerivedFromTheBackgroundColor(t *testing.T) {
	gui := GetDefaultConfig().Gui
	assert.Equal(t, []string{"#626262"}, gui.ThemeForBackground(false, "#1e1e1e").InactiveViewSelectedLineBgColor)
	assert.Equal(t, []string{"#d7d1c1"}, gui.ThemeForBackground(true, "#fdf6e3").InactiveViewSelectedLineBgColor)
}

// If gui.theme had a default for a field that also has a default for the
// terminal's background, the former would always win.
func TestFieldsWithBackgroundDefaultsHaveNoDefaultInGuiTheme(t *testing.T) {
	genericDefaults := setThemeFields(GetDefaultConfig().Gui.Theme)
	for _, field := range setThemeFields(themeDefaults(false, "")) {
		assert.NotContains(t, genericDefaults, field)
	}
	for _, field := range setThemeFields(themeDefaults(true, "")) {
		assert.NotContains(t, genericDefaults, field)
	}
}

// If only one of the backgrounds had a default for a field, the field would
// have no value at all with the other.
func TestDarkAndLightDefaultsSetTheSameFields(t *testing.T) {
	assert.Equal(t, setThemeFields(themeDefaults(false, "")), setThemeFields(themeDefaults(true, "")))
}

func setThemeFields(theme ThemeConfig) []string {
	var fields []string
	themeValue := reflect.ValueOf(theme)
	for i := range themeValue.NumField() {
		if themeValue.Field(i).Len() > 0 {
			fields = append(fields, themeValue.Type().Field(i).Name)
		}
	}
	return fields
}
