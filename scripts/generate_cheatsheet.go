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
    "github.com/jesseduffield/lazygit/pkg/app"
    "github.com/jesseduffield/lazygit/pkg/config"
    "log"
    "os"
    "strings"
)

func writeString(file *os.File, str string) {
    _, err := file.WriteString(str)
    if err != nil {
        log.Fatal(err)
    }
}

func getTitle(mApp *app.App ,viewName string) string {
    viewTitle := strings.Title(viewName) + "Title"
    translatedTitle := mApp.Tr.SLocalize(viewTitle)
    formattedTitle := fmt.Sprintf("\n## %s\n\n", translatedTitle)
    return formattedTitle
}

func main() {
    mConfig, _ := config.NewAppConfig("", "", "", "", "", new(bool))
    mApp, _ := app.Setup(mConfig)
    lang := mApp.Tr.GetLanguage()
    file, _ := os.Create("Keybindings_" + lang + ".md")
    current := ""

    writeString(file, fmt.Sprintf("# Lazygit %s\n", mApp.Tr.SLocalize("menu")))
    writeString(file, getTitle(mApp, "global"))

    writeString(file, "<pre>\n")

    for _, binding := range mApp.Gui.GetKeybindings() {
        if binding.Description == "" {
            continue
        }

        if binding.ViewName != current {
            current = binding.ViewName
            writeString(file, "</pre>\n")
            writeString(file, getTitle(mApp, current))
            writeString(file, "<pre>\n")
        }

        info := fmt.Sprintf("  <kbd>%s</kbd>: %s\n", binding.GetKey(), binding.Description)
        writeString(file, info)
    }

    writeString(file, "</pre>\n")
}
