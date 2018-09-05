// run:
//   LANG=en go run generate_cheatsheet.go
// to generate Keybindings_en.md file in current directory
// change LANG to generate cheatsheet in different language (if supported)

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func main() {
	appConfig, _ := config.NewAppConfig("", "", "", "", "", new(bool))
	a, _ := app.NewApp(appConfig)
	lang := a.Tr.GetLanguage()
	name := "Keybindings_" + lang + ".md"
	bindings := a.Gui.GetKeybindings()
	padWidth := a.Gui.GetMaxKeyLength(bindings)
	file, _ := os.Create(name)
	current := "v"
	content := ""
	title := ""

	file.WriteString("# Lazygit " + a.Tr.SLocalize("menu"))

	for _, binding := range bindings {
		if key := a.Gui.GetKey(binding); key != "" && binding.Description != "" {
			if binding.ViewName != current {
				current = binding.ViewName
				if current == "" {
					title = a.Tr.SLocalize("GlobalTitle")
				} else {
					title = a.Tr.SLocalize(strings.Title(current) + "Title")
				}
				content = fmt.Sprintf("</pre>\n\n## %s\n<pre>\n", title)
				file.WriteString(content)
			}
			content = fmt.Sprintf("\t<kbd>%s</kbd>%s  %s\n", key, strings.TrimPrefix(utils.WithPadding(key, padWidth), key), binding.Description)
			file.WriteString(content)
		}
	}
}
