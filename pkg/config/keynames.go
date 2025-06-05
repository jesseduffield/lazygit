package config

import (
	"strings"
	"unicode/utf8"

	"github.com/jesseduffield/gocui"
	"github.com/samber/lo"
)

// NOTE: if you make changes to this table, be sure to update
// docs/keybindings/Custom_Keybindings.md as well

var LabelByKey = map[gocui.Key]string{
	gocui.KeyF1:             "<f1>",
	gocui.KeyF2:             "<f2>",
	gocui.KeyF3:             "<f3>",
	gocui.KeyF4:             "<f4>",
	gocui.KeyF5:             "<f5>",
	gocui.KeyF6:             "<f6>",
	gocui.KeyF7:             "<f7>",
	gocui.KeyF8:             "<f8>",
	gocui.KeyF9:             "<f9>",
	gocui.KeyF10:            "<f10>",
	gocui.KeyF11:            "<f11>",
	gocui.KeyF12:            "<f12>",
	gocui.KeyInsert:         "<insert>",
	gocui.KeyDelete:         "<delete>",
	gocui.KeyHome:           "<home>",
	gocui.KeyEnd:            "<end>",
	gocui.KeyPgup:           "<pgup>",
	gocui.KeyPgdn:           "<pgdown>",
	gocui.KeyArrowUp:        "<up>",
	gocui.KeyShiftArrowUp:   "<s-up>",
	gocui.KeyArrowDown:      "<down>",
	gocui.KeyShiftArrowDown: "<s-down>",
	gocui.KeyArrowLeft:      "<left>",
	gocui.KeyArrowRight:     "<right>",
	gocui.KeyTab:            "<tab>", // <c-i>
	gocui.KeyBacktab:        "<backtab>",
	gocui.KeyEnter:          "<enter>", // <c-m>
	gocui.KeyAltEnter:       "<a-enter>",
	gocui.KeyEsc:            "<esc>",       // <c-[>, <c-3>
	gocui.KeyBackspace:      "<backspace>", // <c-h>
	gocui.KeyCtrlSpace:      "<c-space>",   // <c-~>, <c-2>
	gocui.KeyCtrlSlash:      "<c-/>",       // <c-_>
	gocui.KeySpace:          "<space>",
	gocui.KeyCtrlA:          "<c-a>",
	gocui.KeyCtrlB:          "<c-b>",
	gocui.KeyCtrlC:          "<c-c>",
	gocui.KeyCtrlD:          "<c-d>",
	gocui.KeyCtrlE:          "<c-e>",
	gocui.KeyCtrlF:          "<c-f>",
	gocui.KeyCtrlG:          "<c-g>",
	gocui.KeyCtrlJ:          "<c-j>",
	gocui.KeyCtrlK:          "<c-k>",
	gocui.KeyCtrlL:          "<c-l>",
	gocui.KeyCtrlN:          "<c-n>",
	gocui.KeyCtrlO:          "<c-o>",
	gocui.KeyCtrlP:          "<c-p>",
	gocui.KeyCtrlQ:          "<c-q>",
	gocui.KeyCtrlR:          "<c-r>",
	gocui.KeyCtrlS:          "<c-s>",
	gocui.KeyCtrlT:          "<c-t>",
	gocui.KeyCtrlU:          "<c-u>",
	gocui.KeyCtrlV:          "<c-v>",
	gocui.KeyCtrlW:          "<c-w>",
	gocui.KeyCtrlX:          "<c-x>",
	gocui.KeyCtrlY:          "<c-y>",
	gocui.KeyCtrlZ:          "<c-z>",
	gocui.KeyCtrl4:          "<c-4>", // <c-\>
	gocui.KeyCtrl5:          "<c-5>", // <c-]>
	gocui.KeyCtrl6:          "<c-6>",
	gocui.KeyCtrl8:          "<c-8>",
	gocui.MouseWheelUp:      "mouse wheel up",
	gocui.MouseWheelDown:    "mouse wheel down",
}

var KeyByLabel = lo.Invert(LabelByKey)

func isValidKeybindingKey(key string) bool {
	runeCount := utf8.RuneCountInString(key)
	if key == "<disabled>" {
		return true
	}

	if runeCount > 1 {
		_, ok := KeyByLabel[strings.ToLower(key)]
		return ok
	}

	return true
}
