package styles

import (
	"sort"

	"github.com/alecthomas/chroma"
)

// Registry of Styles.
var Registry = map[string]*chroma.Style{}

// Fallback style. Reassign to change the default fallback style.
var Fallback = SwapOff

// Register a chroma.Style.
func Register(style *chroma.Style) *chroma.Style {
	Registry[style.Name] = style
	return style
}

// Names of all available styles.
func Names() []string {
	out := []string{}
	for name := range Registry {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// Get named style, or Fallback.
func Get(name string) *chroma.Style {
	if style, ok := Registry[name]; ok {
		return style
	}
	return Fallback
}
