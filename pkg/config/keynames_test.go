package config

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/stretchr/testify/assert"
)

func TestKeyFromLabel(t *testing.T) {
	scenarios := []struct {
		name        string
		label       string
		expectedKey gocui.Key
		expectedOk  bool
	}{
		// Empty / disabled
		{
			name:        "empty string returns unset key",
			label:       "",
			expectedKey: gocui.Key{},
			expectedOk:  true,
		},
		{
			name:        "<disabled> returns unset key",
			label:       "<disabled>",
			expectedKey: gocui.Key{},
			expectedOk:  true,
		},

		// Plain runes (unwrapped)
		{
			name:        "single lowercase letter",
			label:       "a",
			expectedKey: gocui.NewKeyStrMod("a", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "single uppercase letter",
			label:       "A",
			expectedKey: gocui.NewKeyStrMod("A", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "single digit",
			label:       "5",
			expectedKey: gocui.NewKeyStrMod("5", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "punctuation rune",
			label:       "?",
			expectedKey: gocui.NewKeyStrMod("?", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "multibyte rune",
			label:       "ñ",
			expectedKey: gocui.NewKeyStrMod("ñ", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "bare dash is treated as a rune",
			label:       "-",
			expectedKey: gocui.NewKeyRune('-'),
			expectedOk:  true,
		},

		// Special key names (no modifiers, no brackets — though these are
		// always wrapped in brackets in real configs, KeyFromLabel accepts
		// the unwrapped form too)
		{
			name:        "function key",
			label:       "f1",
			expectedKey: gocui.NewKey(gocui.KeyF1, "", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "function key wrapped in brackets",
			label:       "<f12>",
			expectedKey: gocui.NewKey(gocui.KeyF12, "", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "arrow key",
			label:       "<up>",
			expectedKey: gocui.NewKey(gocui.KeyArrowUp, "", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "tab",
			label:       "<tab>",
			expectedKey: gocui.NewKey(gocui.KeyTab, "", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "enter",
			label:       "<enter>",
			expectedKey: gocui.NewKey(gocui.KeyEnter, "", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "esc",
			label:       "<esc>",
			expectedKey: gocui.NewKey(gocui.KeyEsc, "", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "backspace",
			label:       "<backspace>",
			expectedKey: gocui.NewKey(gocui.KeyBackspace, "", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "pgup",
			label:       "<pgup>",
			expectedKey: gocui.NewKey(gocui.KeyPgup, "", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "pgdown",
			label:       "<pgdown>",
			expectedKey: gocui.NewKey(gocui.KeyPgdn, "", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "mouse wheel up",
			label:       "<mouse wheel up>",
			expectedKey: gocui.NewKey(gocui.MouseWheelUp, "", gocui.ModNone),
			expectedOk:  true,
		},

		// Space
		{
			name:        "space keyword maps to space rune",
			label:       "<space>",
			expectedKey: gocui.NewKeyStrMod(" ", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "space keyword without brackets",
			label:       "space",
			expectedKey: gocui.NewKeyStrMod(" ", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "ctrl+space",
			label:       "<c-space>",
			expectedKey: gocui.NewKeyStrMod(" ", gocui.ModCtrl),
			expectedOk:  true,
		},

		// Minus
		{
			name:        "minus keyword maps to dash rune",
			label:       "<minus>",
			expectedKey: gocui.NewKeyStrMod("-", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "ctrl+minus via keyword",
			label:       "<c-minus>",
			expectedKey: gocui.NewKeyStrMod("-", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "ctrl+minus via lenient dash form",
			label:       "<c-->",
			expectedKey: gocui.NewKeyStrMod("-", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "alt+ctrl+minus via lenient dash form",
			label:       "<a-c-->",
			expectedKey: gocui.NewKeyStrMod("-", gocui.ModAlt|gocui.ModCtrl),
			expectedOk:  true,
		},

		// Plus
		{
			name:        "plus keyword maps to plus rune",
			label:       "<plus>",
			expectedKey: gocui.NewKeyStrMod("+", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "ctrl+plus via keyword",
			label:       "<c-plus>",
			expectedKey: gocui.NewKeyStrMod("+", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "ctrl+plus via long keyword and plus separator",
			label:       "<ctrl+plus>",
			expectedKey: gocui.NewKeyStrMod("+", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "alt+shift+plus via keyword",
			label:       "<a-s-plus>",
			expectedKey: gocui.NewKeyStrMod("+", gocui.ModAlt|gocui.ModShift),
			expectedOk:  true,
		},
		{
			name:        "shift alone on plus is rejected",
			label:       "<s-plus>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},

		// Modifiers with runes
		{
			name:        "ctrl+letter",
			label:       "<c-a>",
			expectedKey: gocui.NewKeyStrMod("a", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "alt+letter",
			label:       "<a-x>",
			expectedKey: gocui.NewKeyStrMod("x", gocui.ModAlt),
			expectedOk:  true,
		},
		{
			name:        "meta+letter",
			label:       "<m-z>",
			expectedKey: gocui.NewKeyStrMod("z", gocui.ModMeta),
			expectedOk:  true,
		},

		// Long modifier names are accepted as synonyms for the short forms.
		{
			name:        "ctrl long form",
			label:       "<ctrl-a>",
			expectedKey: gocui.NewKeyStrMod("a", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "alt long form",
			label:       "<alt-x>",
			expectedKey: gocui.NewKeyStrMod("x", gocui.ModAlt),
			expectedOk:  true,
		},
		{
			name:        "meta long form",
			label:       "<meta-z>",
			expectedKey: gocui.NewKeyStrMod("z", gocui.ModMeta),
			expectedOk:  true,
		},
		{
			name:        "shift long form combined with ctrl",
			label:       "<shift-ctrl-a>",
			expectedKey: gocui.NewKeyStrMod("a", gocui.ModShift|gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "long forms work with special keys",
			label:       "<ctrl-up>",
			expectedKey: gocui.NewKey(gocui.KeyArrowUp, "", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "short and long forms can be mixed",
			label:       "<ctrl-s-up>",
			expectedKey: gocui.NewKey(gocui.KeyArrowUp, "", gocui.ModCtrl|gocui.ModShift),
			expectedOk:  true,
		},
		{
			name:        "duplicate via mixed short and long form is rejected",
			label:       "<shift-s-a>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "unknown long modifier is rejected",
			label:       "<control-a>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},

		// Plus is accepted as an alternative modifier separator.
		{
			name:        "plus separator with short form",
			label:       "<c+a>",
			expectedKey: gocui.NewKeyStrMod("a", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "plus separator with long form",
			label:       "<ctrl+alt+a>",
			expectedKey: gocui.NewKeyStrMod("a", gocui.ModCtrl|gocui.ModAlt),
			expectedOk:  true,
		},
		{
			name:        "plus separator with special key",
			label:       "<ctrl+up>",
			expectedKey: gocui.NewKey(gocui.KeyArrowUp, "", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "mixed plus and dash separators",
			label:       "<ctrl+shift-up>",
			expectedKey: gocui.NewKey(gocui.KeyArrowUp, "", gocui.ModCtrl|gocui.ModShift),
			expectedOk:  true,
		},
		{
			name:        "duplicate detection works across separators",
			label:       "<c+c-a>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "ctrl+plus rune via plus separator",
			label:       "<c++>",
			expectedKey: gocui.NewKeyStrMod("+", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "ctrl+dash rune via plus separator",
			label:       "<c+->",
			expectedKey: gocui.NewKeyStrMod("-", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "ctrl+plus rune via dash separator",
			label:       "<c-+>",
			expectedKey: gocui.NewKeyStrMod("+", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "bare plus rune",
			label:       "+",
			expectedKey: gocui.NewKeyStrMod("+", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "bare plus wrapped in brackets",
			label:       "<+>",
			expectedKey: gocui.NewKeyStrMod("+", gocui.ModNone),
			expectedOk:  true,
		},

		// Shift-on-rune is rejected: terminals fold shift into the rune
		// itself, so the binding could never fire. Combined with other
		// modifiers it's allowed (the terminal can't fold it then).
		{
			name:        "shift alone on a letter is rejected",
			label:       "<s-q>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "shift alone on uppercase letter is rejected",
			label:       "<s-A>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "shift alone on minus is rejected",
			label:       "<s-minus>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "shift on space is allowed (rune does not change)",
			label:       "<s-space>",
			expectedKey: gocui.NewKeyStrMod(" ", gocui.ModShift),
			expectedOk:  true,
		},
		{
			name:        "shift combined with ctrl on a letter is allowed",
			label:       "<c-s-x>",
			expectedKey: gocui.NewKeyStrMod("x", gocui.ModCtrl|gocui.ModShift),
			expectedOk:  true,
		},
		{
			name:        "shift combined with alt on minus is allowed",
			label:       "<a-s-minus>",
			expectedKey: gocui.NewKeyStrMod("-", gocui.ModAlt|gocui.ModShift),
			expectedOk:  true,
		},

		// Uppercase ASCII letter with a modifier is rejected: ctrl+letter
		// always arrives with a lowercase rune (control codes have no case
		// distinction), and CSI-u reports the unshifted codepoint with
		// shift as a separate modifier.
		{
			name:        "ctrl+uppercase letter is rejected",
			label:       "<c-A>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "alt+uppercase letter is rejected",
			label:       "<a-A>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "meta+uppercase letter is rejected",
			label:       "<m-A>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "combined modifier on uppercase letter is rejected",
			label:       "<c-a-A>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "bare uppercase letter is allowed",
			label:       "A",
			expectedKey: gocui.NewKeyStrMod("A", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "modifier on digit is allowed",
			label:       "<c-1>",
			expectedKey: gocui.NewKeyStrMod("1", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "modifier on non-ASCII uppercase letter is allowed",
			label:       "<a-Ñ>",
			expectedKey: gocui.NewKeyStrMod("Ñ", gocui.ModAlt),
			expectedOk:  true,
		},

		// Modifiers with special keys
		{
			name:        "ctrl+enter",
			label:       "<c-enter>",
			expectedKey: gocui.NewKey(gocui.KeyEnter, "", gocui.ModCtrl),
			expectedOk:  true,
		},
		{
			name:        "alt+up",
			label:       "<a-up>",
			expectedKey: gocui.NewKey(gocui.KeyArrowUp, "", gocui.ModAlt),
			expectedOk:  true,
		},
		{
			name:        "shift+f1",
			label:       "<s-f1>",
			expectedKey: gocui.NewKey(gocui.KeyF1, "", gocui.ModShift),
			expectedOk:  true,
		},
		{
			name:        "meta+enter",
			label:       "<m-enter>",
			expectedKey: gocui.NewKey(gocui.KeyEnter, "", gocui.ModMeta),
			expectedOk:  true,
		},

		// Combined modifiers
		{
			name:        "ctrl+alt+letter",
			label:       "<c-a-x>",
			expectedKey: gocui.NewKeyStrMod("x", gocui.ModCtrl|gocui.ModAlt),
			expectedOk:  true,
		},
		{
			name:        "all four modifiers on a letter",
			label:       "<s-c-a-m-x>",
			expectedKey: gocui.NewKeyStrMod("x", gocui.ModShift|gocui.ModCtrl|gocui.ModAlt|gocui.ModMeta),
			expectedOk:  true,
		},
		{
			name:        "ctrl+shift+arrow key",
			label:       "<c-s-up>",
			expectedKey: gocui.NewKey(gocui.KeyArrowUp, "", gocui.ModCtrl|gocui.ModShift),
			expectedOk:  true,
		},

		// Bracket handling
		{
			name:        "single rune wrapped in brackets is unwrapped",
			label:       "<a>",
			expectedKey: gocui.NewKeyStrMod("a", gocui.ModNone),
			expectedOk:  true,
		},
		{
			name:        "dash wrapped in brackets",
			label:       "<->",
			expectedKey: gocui.NewKeyRune('-'),
			expectedOk:  true,
		},

		// Invalid inputs
		{
			name:        "unknown special key name",
			label:       "<nope>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "unknown modifier letter",
			label:       "<x-a>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "uppercase modifier is not accepted",
			label:       "<C-a>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "duplicate ctrl modifier",
			label:       "<c-c-a>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "duplicate shift modifier",
			label:       "<s-s-a>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "duplicate alt modifier",
			label:       "<a-a-x>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "duplicate meta modifier",
			label:       "<m-m-x>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "trailing modifier with no key",
			label:       "<c->",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "multi-character non-special label",
			label:       "ab",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "empty brackets",
			label:       "<>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
		{
			name:        "modifier on unknown key name",
			label:       "<c-nope>",
			expectedKey: gocui.Key{},
			expectedOk:  false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			key, ok := KeyFromLabel(s.label)
			assert.Equal(t, s.expectedOk, ok)
			assert.Equal(t, s.expectedKey, key)
		})
	}
}

func TestLabelForKey(t *testing.T) {
	scenarios := []struct {
		name     string
		key      gocui.Key
		expected string
	}{
		// Unset
		{"unset key produces empty string", gocui.Key{}, ""},

		// Plain runes — single-character output, no brackets
		{"lowercase letter", gocui.NewKeyStrMod("a", gocui.ModNone), "a"},
		{"uppercase letter", gocui.NewKeyStrMod("A", gocui.ModNone), "A"},
		{"digit", gocui.NewKeyStrMod("5", gocui.ModNone), "5"},
		{"punctuation", gocui.NewKeyStrMod("?", gocui.ModNone), "?"},
		{"slash", gocui.NewKeyStrMod("/", gocui.ModNone), "/"},
		{"multibyte rune", gocui.NewKeyStrMod("ñ", gocui.ModNone), "ñ"},

		// Space and dash — special-cased rune output
		{"plain dash uses literal", gocui.NewKeyStrMod("-", gocui.ModNone), "-"},
		{"plain space uses keyword", gocui.NewKeyStrMod(" ", gocui.ModNone), "<space>"},
		{"ctrl+dash uses minus keyword", gocui.NewKeyStrMod("-", gocui.ModCtrl), "<ctrl+minus>"},
		{"alt+dash uses minus keyword", gocui.NewKeyStrMod("-", gocui.ModAlt), "<alt+minus>"},
		{"plain plus uses literal", gocui.NewKeyStrMod("+", gocui.ModNone), "+"},
		{"ctrl+plus uses plus keyword", gocui.NewKeyStrMod("+", gocui.ModCtrl), "<ctrl+plus>"},
		{"alt+plus uses plus keyword", gocui.NewKeyStrMod("+", gocui.ModAlt), "<alt+plus>"},
		{"ctrl+space", gocui.NewKeyStrMod(" ", gocui.ModCtrl), "<ctrl+space>"},

		// Single modifier on a rune
		{"ctrl+letter", gocui.NewKeyStrMod("a", gocui.ModCtrl), "<ctrl+a>"},
		{"alt+letter", gocui.NewKeyStrMod("x", gocui.ModAlt), "<alt+x>"},
		{"meta+letter", gocui.NewKeyStrMod("z", gocui.ModMeta), "<meta+z>"},
		{"shift+space", gocui.NewKeyStrMod(" ", gocui.ModShift), "<shift+space>"},

		// Modifier ordering — canonical output is ctrl+, alt+, shift+, meta+
		{"ctrl+alt orders ctrl before alt", gocui.NewKeyStrMod("x", gocui.ModCtrl|gocui.ModAlt), "<ctrl+alt+x>"},
		{"shift+ctrl orders ctrl before shift", gocui.NewKeyStrMod("x", gocui.ModShift|gocui.ModCtrl), "<ctrl+shift+x>"},
		{"meta+shift orders shift before meta", gocui.NewKeyStrMod("x", gocui.ModMeta|gocui.ModShift), "<shift+meta+x>"},
		{
			"all four modifiers ordered ctrl+alt+shift+meta",
			gocui.NewKeyStrMod("x", gocui.ModCtrl|gocui.ModAlt|gocui.ModShift|gocui.ModMeta),
			"<ctrl+alt+shift+meta+x>",
		},

		// Special keys (always wrapped, even unmodified)
		{"f1", gocui.NewKey(gocui.KeyF1, "", gocui.ModNone), "<f1>"},
		{"f12", gocui.NewKey(gocui.KeyF12, "", gocui.ModNone), "<f12>"},
		{"insert", gocui.NewKey(gocui.KeyInsert, "", gocui.ModNone), "<insert>"},
		{"delete", gocui.NewKey(gocui.KeyDelete, "", gocui.ModNone), "<delete>"},
		{"home", gocui.NewKey(gocui.KeyHome, "", gocui.ModNone), "<home>"},
		{"end", gocui.NewKey(gocui.KeyEnd, "", gocui.ModNone), "<end>"},
		{"pgup", gocui.NewKey(gocui.KeyPgup, "", gocui.ModNone), "<pgup>"},
		{"pgdown", gocui.NewKey(gocui.KeyPgdn, "", gocui.ModNone), "<pgdown>"},
		{"arrow up", gocui.NewKey(gocui.KeyArrowUp, "", gocui.ModNone), "<up>"},
		{"arrow down", gocui.NewKey(gocui.KeyArrowDown, "", gocui.ModNone), "<down>"},
		{"arrow left", gocui.NewKey(gocui.KeyArrowLeft, "", gocui.ModNone), "<left>"},
		{"arrow right", gocui.NewKey(gocui.KeyArrowRight, "", gocui.ModNone), "<right>"},
		{"tab", gocui.NewKey(gocui.KeyTab, "", gocui.ModNone), "<tab>"},
		{"backtab", gocui.NewKey(gocui.KeyBacktab, "", gocui.ModNone), "<backtab>"},
		{"enter", gocui.NewKey(gocui.KeyEnter, "", gocui.ModNone), "<enter>"},
		{"esc", gocui.NewKey(gocui.KeyEsc, "", gocui.ModNone), "<esc>"},
		{"backspace", gocui.NewKey(gocui.KeyBackspace, "", gocui.ModNone), "<backspace>"},
		{"mouse wheel up", gocui.NewKey(gocui.MouseWheelUp, "", gocui.ModNone), "<mouse wheel up>"},
		{"mouse wheel down", gocui.NewKey(gocui.MouseWheelDown, "", gocui.ModNone), "<mouse wheel down>"},

		// Modifiers on special keys
		{"shift+f1", gocui.NewKey(gocui.KeyF1, "", gocui.ModShift), "<shift+f1>"},
		{"alt+up", gocui.NewKey(gocui.KeyArrowUp, "", gocui.ModAlt), "<alt+up>"},
		{"meta+enter", gocui.NewKey(gocui.KeyEnter, "", gocui.ModMeta), "<meta+enter>"},
		{"ctrl+shift+up", gocui.NewKey(gocui.KeyArrowUp, "", gocui.ModCtrl|gocui.ModShift), "<ctrl+shift+up>"},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, LabelForKey(s.key))
		})
	}
}

// Round-trip: every label produced by LabelForKey should parse back to the
// same key via KeyFromLabel.
func TestKeyFromLabel_RoundTripFromLabelForKey(t *testing.T) {
	scenarios := []struct {
		name string
		key  gocui.Key
	}{
		{"unset key", gocui.Key{}},
		{"plain letter", gocui.NewKeyStrMod("a", gocui.ModNone)},
		{"plain digit", gocui.NewKeyStrMod("7", gocui.ModNone)},
		{"space", gocui.NewKeyStrMod(" ", gocui.ModNone)},
		{"ctrl+letter", gocui.NewKeyStrMod("a", gocui.ModCtrl)},
		{"alt+letter", gocui.NewKeyStrMod("x", gocui.ModAlt)},
		{"meta+letter", gocui.NewKeyStrMod("z", gocui.ModMeta)},
		{"shift+space", gocui.NewKeyStrMod(" ", gocui.ModShift)},
		{"ctrl+shift+letter", gocui.NewKeyStrMod("x", gocui.ModCtrl|gocui.ModShift)},
		{"ctrl+alt+letter", gocui.NewKeyStrMod("x", gocui.ModCtrl|gocui.ModAlt)},
		{"f1", gocui.NewKey(gocui.KeyF1, "", gocui.ModNone)},
		{"shift+f1", gocui.NewKey(gocui.KeyF1, "", gocui.ModShift)},
		{"alt+up", gocui.NewKey(gocui.KeyArrowUp, "", gocui.ModAlt)},
		{"meta+enter", gocui.NewKey(gocui.KeyEnter, "", gocui.ModMeta)},
		{"esc", gocui.NewKey(gocui.KeyEsc, "", gocui.ModNone)},
		{"mouse wheel up", gocui.NewKey(gocui.MouseWheelUp, "", gocui.ModNone)},
		{"ctrl+space", gocui.NewKeyStrMod(" ", gocui.ModCtrl)},
		{"plain dash", gocui.NewKeyStrMod("-", gocui.ModNone)},
		{"ctrl+dash", gocui.NewKeyStrMod("-", gocui.ModCtrl)},
		{"alt+shift+dash", gocui.NewKeyStrMod("-", gocui.ModAlt|gocui.ModShift)},
		{"plain plus", gocui.NewKeyStrMod("+", gocui.ModNone)},
		{"ctrl+plus", gocui.NewKeyStrMod("+", gocui.ModCtrl)},
		{"alt+shift+plus", gocui.NewKeyStrMod("+", gocui.ModAlt|gocui.ModShift)},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			label := LabelForKey(s.key)
			parsed, ok := KeyFromLabel(label)
			assert.True(t, ok, "expected label %q to parse", label)
			assert.Equal(t, s.key, parsed)
		})
	}
}
