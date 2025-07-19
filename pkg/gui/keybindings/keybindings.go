package keybindings

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func Label(name string) string {
	return LabelFromKey(GetKey(name))
}

func LabelFromKey(key types.Key) string {
	if key == nil {
		return ""
	}

	keyInt := 0

	switch key := key.(type) {
	case rune:
		keyInt = int(key)
	case gocui.Key:
		value, ok := config.LabelByKey[key]
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
		binding, ok := config.KeyByLabel[strings.ToLower(key)]
		if !ok {
			log.Fatalf("Unrecognized key %s for keybinding. For permitted values see %s", strings.ToLower(key), constants.Links.Docs.CustomKeybindings)
		}
		return binding
	} else if runeCount == 1 {
		return []rune(key)[0]
	}
	return nil
}
