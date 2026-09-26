package config

import (
	"fmt"
	"maps"
	"math"
	"reflect"
	"slices"
	"strconv"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// ThemeForBackground returns gui.theme with the overrides for a dark or a light
// background applied, and with the defaults for that background in the fields
// that neither of them sets. backgroundColor is the terminal's background color
// as #rrggbb, or empty if we don't know it.
func (c *GuiConfig) ThemeForBackground(lightBackground bool, backgroundColor string) ThemeConfig {
	override := lo.Ternary(lightBackground, c.LightTheme, c.DarkTheme)
	return mergeThemes(override, c.Theme, themeDefaults(lightBackground, backgroundColor))
}

// themeDefaults returns the defaults of the theme fields whose default depends
// on the terminal's background. GetDefaultConfig leaves these fields empty in
// gui.theme, so that a value there comes from the user, and wins over these.
//
// They are derived from the background color, so that they keep the same
// distance from it however dark or light it is. If we don't know the
// background color, we assume black or white.
func themeDefaults(lightBackground bool, backgroundColor string) ThemeConfig {
	if len(backgroundColor) != 7 || !utils.IsValidHexValue(backgroundColor) {
		backgroundColor = lo.Ternary(lightBackground, "#ffffff", "#000000")
	}

	if lightBackground {
		return ThemeConfig{
			InactiveViewSelectedLineBgColor: []string{mixHexColors(backgroundColor, "#000000", 0.15)},
		}
	}

	return ThemeConfig{
		InactiveViewSelectedLineBgColor: []string{mixHexColors(backgroundColor, "#ffffff", 0.3)},
	}
}

// mixHexColors mixes two colors given as #rrggbb. amount is how much of b to
// take: 0 gives a, 1 gives b.
func mixHexColors(a, b string, amount float64) string {
	channel := func(color string, i int) float64 {
		value, _ := strconv.ParseUint(color[1+2*i:3+2*i], 16, 8)
		return float64(value)
	}

	result := "#"
	for i := range 3 {
		mixed := channel(a, i) + (channel(b, i)-channel(a, i))*amount
		result += fmt.Sprintf("%02x", int(math.Round(mixed)))
	}
	return result
}

// mergeThemes takes each field from the first of the themes that sets it. For
// maps and color patterns it merges the entries instead, and an entry of an
// earlier theme wins over one with the same key in a later theme.
func mergeThemes(themes ...ThemeConfig) ThemeConfig {
	var result ThemeConfig
	resultValue := reflect.ValueOf(&result).Elem()
	for i := range resultValue.NumField() {
		values := lo.Map(themes, func(theme ThemeConfig, _ int) any {
			return reflect.ValueOf(theme).Field(i).Interface()
		})
		resultValue.Field(i).Set(reflect.ValueOf(mergeThemeField(values)))
	}
	return result
}

func mergeThemeField(values []any) any {
	switch values[0].(type) {
	case []string:
		return lo.FindOrElse(lo.Map(values, func(v any, _ int) []string { return v.([]string) }), nil,
			func(v []string) bool { return len(v) > 0 })
	case map[string]string:
		merged := map[string]string{}
		for _, v := range slices.Backward(values) {
			maps.Copy(merged, v.(map[string]string))
		}
		return merged
	case ColorPatterns:
		return lo.Reduce(values, func(merged ColorPatterns, v any, _ int) ColorPatterns {
			return merged.over(v.(ColorPatterns))
		}, nil)
	default:
		panic(fmt.Sprintf("don't know how to merge a theme field of type %T", values[0]))
	}
}
