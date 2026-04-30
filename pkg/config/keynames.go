package config

import (
	"log"
	"strings"
	"unicode/utf8"

	"github.com/gdamore/tcell/v3"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/samber/lo"
)

// NOTE: if you make changes to this table, be sure to update
// docs/keybindings/Custom_Keybindings.md as well

var labelByKey = map[gocui.KeyName]string{
	gocui.KeyF1:          "f1",
	gocui.KeyF2:          "f2",
	gocui.KeyF3:          "f3",
	gocui.KeyF4:          "f4",
	gocui.KeyF5:          "f5",
	gocui.KeyF6:          "f6",
	gocui.KeyF7:          "f7",
	gocui.KeyF8:          "f8",
	gocui.KeyF9:          "f9",
	gocui.KeyF10:         "f10",
	gocui.KeyF11:         "f11",
	gocui.KeyF12:         "f12",
	gocui.KeyInsert:      "insert",
	gocui.KeyDelete:      "delete",
	gocui.KeyHome:        "home",
	gocui.KeyEnd:         "end",
	gocui.KeyPgup:        "pgup",
	gocui.KeyPgdn:        "pgdown",
	gocui.KeyArrowUp:     "up",
	gocui.KeyArrowDown:   "down",
	gocui.KeyArrowLeft:   "left",
	gocui.KeyArrowRight:  "right",
	gocui.KeyTab:         "tab",
	gocui.KeyBacktab:     "backtab",
	gocui.KeyEnter:       "enter",
	gocui.KeyEsc:         "esc",
	gocui.KeyBackspace:   "backspace",
	gocui.MouseWheelUp:   "mouse wheel up",
	gocui.MouseWheelDown: "mouse wheel down",
}

var keyByLabel = lo.Invert(labelByKey)

func LabelForKey(key gocui.Key) string {
	if !key.IsSet() {
		return ""
	}

	label := ""
	if key.Mod()&gocui.ModCtrl != 0 {
		label += "ctrl+"
	}
	if key.Mod()&gocui.ModAlt != 0 {
		label += "alt+"
	}
	if key.Mod()&gocui.ModShift != 0 {
		label += "shift+"
	}
	if key.Mod()&gocui.ModMeta != 0 {
		label += "meta+"
	}

	if key.KeyName() == gocui.KeyName(tcell.KeyRune) {
		if key.Str() == " " {
			label += "space"
		} else if key.Str() == "-" && key.Mod() != gocui.ModNone {
			label += "minus"
		} else if key.Str() == "+" && key.Mod() != gocui.ModNone {
			label += "plus"
		} else {
			label += key.Str()
		}
	} else {
		value, ok := labelByKey[key.KeyName()]
		if ok {
			label += value
		} else {
			label += "unknown"
		}
	}

	if utf8.RuneCountInString(label) > 1 {
		label = "<" + label + ">"
	}

	return label
}

func KeyFromLabel(label string) (gocui.Key, bool) {
	if label == "" || label == "<disabled>" {
		return gocui.Key{}, true
	}

	if strings.HasPrefix(label, "<") && strings.HasSuffix(label, ">") {
		label = label[1 : len(label)-1]
	}

	mod := gocui.ModNone
	for {
		// A bare "-" or "+" with any (or no) modifiers is a literal rune
		// key; this also covers lenient forms like `<c-->` and `<c++>`,
		// neither of which we emit (we use `<c-minus>` and `<c-+>`).
		if label == "-" || label == "+" {
			return gocui.NewKeyStrMod(label, mod), true
		}

		sepIdx := strings.IndexAny(label, "-+")
		if sepIdx == -1 {
			break
		}
		modStr, remainder := label[:sepIdx], label[sepIdx+1:]

		label = remainder

		switch modStr {
		case "s", "shift":
			if (mod & gocui.ModShift) != 0 {
				return gocui.Key{}, false
			}
			mod |= gocui.ModShift
		case "c", "ctrl":
			if (mod & gocui.ModCtrl) != 0 {
				return gocui.Key{}, false
			}
			mod |= gocui.ModCtrl
		case "a", "alt":
			if (mod & gocui.ModAlt) != 0 {
				return gocui.Key{}, false
			}
			mod |= gocui.ModAlt
		case "m", "meta":
			if (mod & gocui.ModMeta) != 0 {
				return gocui.Key{}, false
			}
			mod |= gocui.ModMeta
		default:
			return gocui.Key{}, false
		}
	}

	if label == "space" {
		return gocui.NewKeyStrMod(" ", mod), true
	}

	if label == "minus" {
		if mod == gocui.ModShift {
			return gocui.Key{}, false
		}
		return gocui.NewKeyStrMod("-", mod), true
	}

	if label == "plus" {
		if mod == gocui.ModShift {
			return gocui.Key{}, false
		}
		return gocui.NewKeyStrMod("+", mod), true
	}

	if keyName, ok := keyByLabel[label]; ok {
		return gocui.NewKey(keyName, "", mod), true
	}

	runeCount := utf8.RuneCountInString(label)
	if runeCount != 1 {
		return gocui.Key{}, false
	}

	// Shift on a bare rune is invalid: terminals fold shift into the rune
	// itself (shift+a arrives as "A"), so the binding could never fire.
	// Space is exempt and handled above; combined with other modifiers,
	// shift is fine because the terminal can't fold it into the rune then.
	if mod == gocui.ModShift {
		return gocui.Key{}, false
	}

	// An ASCII uppercase letter with any modifier is invalid. Ctrl+letter
	// events always arrive with a lowercase rune — control codes have no
	// case distinction (the terminal sends the same byte for ctrl+a and
	// ctrl+A), and CSI-u protocols report the unshifted codepoint with
	// shift as a separate modifier (alt+shift+a → rune='a' mod=Alt|Shift).
	// Users should write <c-s-a> rather than <c-A>.
	if mod != gocui.ModNone && len(label) == 1 && label[0] >= 'A' && label[0] <= 'Z' {
		return gocui.Key{}, false
	}

	return gocui.NewKeyStrMod(label, mod), true
}

func isValidKeybindingKey(key string) bool {
	_, ok := KeyFromLabel(key)
	return ok
}

func GetValidatedKeyBindingKey(label string) gocui.Key {
	key, ok := KeyFromLabel(label)
	if !ok {
		log.Fatalf("Unrecognized key %s, this should have been caught by user config validation", label)
	}

	return key
}
