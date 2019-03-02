// This "script" generates a file called Keybindings_{{.LANG}}.md
// in current working directory.
//
// The content of this generated file is a keybindings cheatsheet.
//
// To generate cheatsheet in english run:
//   LANG=en go run scripts/generate_cheatsheet.go

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui"
)

func writeString(file *os.File, str string) {
	_, err := file.WriteString(str)
	if err != nil {
		log.Fatal(err)
	}
}

func localisedTitle(mApp *app.App, str string) string {
	viewTitle := strings.Title(str) + "Title"
	return mApp.Tr.SLocalize(viewTitle)
}

func formatTitle(title string) string {
	return fmt.Sprintf("\n## %s\n\n", title)
}

func writeBinding(file *os.File, binding *gui.Binding) {
	info := fmt.Sprintf("  <kbd>%s</kbd>: %s\n", binding.GetKey(), binding.Description)
	writeString(file, info)
}

// I should really just build an array of tuples, one thing with a string and the other with a list of bindings, and then build them like that.

func main() {
	mConfig, _ := config.NewAppConfig("", "", "", "", "", true)
	mApp, _ := app.NewApp(mConfig)
	lang := mApp.Tr.GetLanguage()
	file, _ := os.Create("Keybindings_" + lang + ".md")
	current := ""

	writeString(file, fmt.Sprintf("# Lazygit %s\n", mApp.Tr.SLocalize("menu")))
	writeString(file, formatTitle(localisedTitle(mApp, "global")))

	writeString(file, "<pre>\n")

	// TODO: add context-based keybindings
	for _, binding := range mApp.Gui.GetInitialKeybindings() {
		if binding.Description == "" {
			continue
		}

		if binding.ViewName != current {
			current = binding.ViewName
			writeString(file, "</pre>\n")
			writeString(file, formatTitle(localisedTitle(mApp, current)))
			writeString(file, "<pre>\n")
		}

		writeBinding(file, binding)
	}

	writeString(file, "</pre>\n")

	for view, contexts := range mApp.Gui.GetContextMap() {
		for contextName, contextBindings := range contexts {
			translatedView := localisedTitle(mApp, view)
			translatedContextName := localisedTitle(mApp, contextName)
			writeString(file, fmt.Sprintf("\n## %s (%s)\n\n", translatedView, translatedContextName))
			writeString(file, "<pre>\n")
			for _, binding := range contextBindings {
				if binding.Description == "" {
					continue
				}

				writeBinding(file, binding)
			}
			writeString(file, "</pre>\n")
		}
	}
}
