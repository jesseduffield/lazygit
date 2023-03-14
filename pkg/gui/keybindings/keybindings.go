package keybindings

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

var keyMapReversed = map[gocui.Key]string{
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
	gocui.KeyArrowUp:     "▲",
	gocui.KeyArrowDown:   "▼",
	gocui.KeyArrowLeft:   "◀",
	gocui.KeyArrowRight:  "▶",
	gocui.KeyTab:         "tab", // ctrl+i
	gocui.KeyBacktab:     "shift+tab",
	gocui.KeyEnter:       "enter", // ctrl+m
	gocui.KeyAltEnter:    "alt+enter",
	gocui.KeyEsc:         "esc",        // ctrl+[, ctrl+3
	gocui.KeyBackspace:   "backspace",  // ctrl+h
	gocui.KeyCtrlSpace:   "ctrl+space", // ctrl+~, ctrl+2
	gocui.KeyCtrlSlash:   "ctrl+/",     // ctrl+_
	gocui.KeySpace:       "space",
	gocui.KeyCtrlA:       "ctrl+a",
	gocui.KeyCtrlB:       "ctrl+b",
	gocui.KeyCtrlC:       "ctrl+c",
	gocui.KeyCtrlD:       "ctrl+d",
	gocui.KeyCtrlE:       "ctrl+e",
	gocui.KeyCtrlF:       "ctrl+f",
	gocui.KeyCtrlG:       "ctrl+g",
	gocui.KeyCtrlJ:       "ctrl+j",
	gocui.KeyCtrlK:       "ctrl+k",
	gocui.KeyCtrlL:       "ctrl+l",
	gocui.KeyCtrlN:       "ctrl+n",
	gocui.KeyCtrlO:       "ctrl+o",
	gocui.KeyCtrlP:       "ctrl+p",
	gocui.KeyCtrlQ:       "ctrl+q",
	gocui.KeyCtrlR:       "ctrl+r",
	gocui.KeyCtrlS:       "ctrl+s",
	gocui.KeyCtrlT:       "ctrl+t",
	gocui.KeyCtrlU:       "ctrl+u",
	gocui.KeyCtrlV:       "ctrl+v",
	gocui.KeyCtrlW:       "ctrl+w",
	gocui.KeyCtrlX:       "ctrl+x",
	gocui.KeyCtrlY:       "ctrl+y",
	gocui.KeyCtrlZ:       "ctrl+z",
	gocui.KeyCtrl4:       "ctrl+4", // ctrl+\
	gocui.KeyCtrl5:       "ctrl+5", // ctrl+]
	gocui.KeyCtrl6:       "ctrl+6",
	gocui.KeyCtrl8:       "ctrl+8",
	gocui.MouseWheelUp:   "mouse wheel ▲",
	gocui.MouseWheelDown: "mouse wheel ▼",
}

var keyMap = map[string]types.Key{
	"<c-a>":       gocui.KeyCtrlA,
	"<c-b>":       gocui.KeyCtrlB,
	"<c-c>":       gocui.KeyCtrlC,
	"<c-d>":       gocui.KeyCtrlD,
	"<c-e>":       gocui.KeyCtrlE,
	"<c-f>":       gocui.KeyCtrlF,
	"<c-g>":       gocui.KeyCtrlG,
	"<c-h>":       gocui.KeyCtrlH,
	"<c-i>":       gocui.KeyCtrlI,
	"<c-j>":       gocui.KeyCtrlJ,
	"<c-k>":       gocui.KeyCtrlK,
	"<c-l>":       gocui.KeyCtrlL,
	"<c-m>":       gocui.KeyCtrlM,
	"<c-n>":       gocui.KeyCtrlN,
	"<c-o>":       gocui.KeyCtrlO,
	"<c-p>":       gocui.KeyCtrlP,
	"<c-q>":       gocui.KeyCtrlQ,
	"<c-r>":       gocui.KeyCtrlR,
	"<c-s>":       gocui.KeyCtrlS,
	"<c-t>":       gocui.KeyCtrlT,
	"<c-u>":       gocui.KeyCtrlU,
	"<c-v>":       gocui.KeyCtrlV,
	"<c-w>":       gocui.KeyCtrlW,
	"<c-x>":       gocui.KeyCtrlX,
	"<c-y>":       gocui.KeyCtrlY,
	"<c-z>":       gocui.KeyCtrlZ,
	"<c-~>":       gocui.KeyCtrlTilde,
	"<c-2>":       gocui.KeyCtrl2,
	"<c-3>":       gocui.KeyCtrl3,
	"<c-4>":       gocui.KeyCtrl4,
	"<c-5>":       gocui.KeyCtrl5,
	"<c-6>":       gocui.KeyCtrl6,
	"<c-7>":       gocui.KeyCtrl7,
	"<c-8>":       gocui.KeyCtrl8,
	"<c-space>":   gocui.KeyCtrlSpace,
	"<c-\\>":      gocui.KeyCtrlBackslash,
	"<c-[>":       gocui.KeyCtrlLsqBracket,
	"<c-]>":       gocui.KeyCtrlRsqBracket,
	"<c-/>":       gocui.KeyCtrlSlash,
	"<c-_>":       gocui.KeyCtrlUnderscore,
	"<backspace>": gocui.KeyBackspace,
	"<tab>":       gocui.KeyTab,
	"<backtab>":   gocui.KeyBacktab,
	"<enter>":     gocui.KeyEnter,
	"<a-enter>":   gocui.KeyAltEnter,
	"<esc>":       gocui.KeyEsc,
	"<space>":     gocui.KeySpace,
	"<f1>":        gocui.KeyF1,
	"<f2>":        gocui.KeyF2,
	"<f3>":        gocui.KeyF3,
	"<f4>":        gocui.KeyF4,
	"<f5>":        gocui.KeyF5,
	"<f6>":        gocui.KeyF6,
	"<f7>":        gocui.KeyF7,
	"<f8>":        gocui.KeyF8,
	"<f9>":        gocui.KeyF9,
	"<f10>":       gocui.KeyF10,
	"<f11>":       gocui.KeyF11,
	"<f12>":       gocui.KeyF12,
	"<insert>":    gocui.KeyInsert,
	"<delete>":    gocui.KeyDelete,
	"<home>":      gocui.KeyHome,
	"<end>":       gocui.KeyEnd,
	"<pgup>":      gocui.KeyPgup,
	"<pgdown>":    gocui.KeyPgdn,
	"<up>":        gocui.KeyArrowUp,
	"<down>":      gocui.KeyArrowDown,
	"<left>":      gocui.KeyArrowLeft,
	"<right>":     gocui.KeyArrowRight,
}

func Label(name string) string {
	return LabelFromKey(GetKey(name))
}

func LabelFromKey(key types.Key) string {
	keyInt := 0

	switch key := key.(type) {
	case rune:
		keyInt = int(key)
	case gocui.Key:
		value, ok := keyMapReversed[key]
		if ok {
			return value
		}
		keyInt = int(key)
	}

	return fmt.Sprintf("%c", keyInt)
}

func GetKey(key string) types.Key {
	runeCount := utf8.RuneCountInString(key)
	if runeCount > 1 {
		binding := keyMap[strings.ToLower(key)]
		if binding == nil {
			log.Fatalf("Unrecognized key %s for keybinding. For permitted values see %s", strings.ToLower(key), constants.Links.Docs.CustomKeybindings)
		} else {
			return binding
		}
	} else if runeCount == 1 {
		return []rune(key)[0]
	}
	return nil
}
