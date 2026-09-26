package config

import (
	"fmt"
	"maps"
	"reflect"
	"slices"

	"github.com/samber/lo"
)

// ThemeForBackground returns gui.theme with the overrides for a dark or a light
// background applied.
func (c *GuiConfig) ThemeForBackground(lightBackground bool) ThemeConfig {
	override := lo.Ternary(lightBackground, c.LightTheme, c.DarkTheme)
	return mergeThemes(override, c.Theme)
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
