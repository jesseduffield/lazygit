package gui

import (
	"os"
	"regexp"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

func firstTruthyString(strs ...string) string {
	for _, str := range strs {
		if str != "" {
			return str
		}
	}
	return ""
}

func canUseUnicode() bool {
	re := regexp.MustCompile(`(?i)utf-8`)
	return re.MatchString(
		firstTruthyString(os.Getenv("LC_ALL"), os.Getenv("LC_CTYPE"), os.Getenv("LANG")),
	)
}

func (gui *Gui) prepareEncodings() {
	if canUseUnicode() {
		gui.encodedStrings = &utils.EncodedStrings{
			UpArrow:    "↑",
			DownArrow:  "↓",
			LeftArrow:  "←",
			RightArrow: "→",
		}
	} else {
		gui.g.ASCII = true
		gui.encodedStrings = &utils.EncodedStrings{
			UpArrow:    "^",
			DownArrow:  "v",
			LeftArrow:  "<-",
			RightArrow: "->",
		}
	}
}
