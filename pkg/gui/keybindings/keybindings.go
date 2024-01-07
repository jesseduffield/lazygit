package keybindings

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

var labelByKey = map[gocui.Key]string{
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

var keyByLabel = lo.Invert(labelByKey)

func Label(name string) string {
	return LabelFromKey(GetKey(name))
}

func LabelFromKey(key types.Key) string {
	keyInt := 0

	switch key := key.(type) {
	case rune:
		keyInt = int(key)
	case gocui.Key:
		value, ok := labelByKey[key]
		if ok {
			return value
		}
		keyInt = int(key)
	}

	return fmt.Sprintf("%c", keyInt)
}

func GetKey(key string) types.Key {
	runeCount := utf8.RuneCountInString(key)
	if key == "<disabled>" {
		return nil
	} else if runeCount > 1 {
		binding, ok := keyByLabel[strings.ToLower(key)]
		if !ok {
			log.Fatalf("Unrecognized key %s for keybinding. For permitted values see %s", strings.ToLower(key), constants.Links.Docs.CustomKeybindings)
		} else {
			return binding
		}
	} else if runeCount == 1 {
		return []rune(key)[0]
	}
	return nil
}
