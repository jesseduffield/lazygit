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
	"sort"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui"
)

type bindingSection struct {
	title    string
	bindings []*gui.Binding
}

func main() {
	langs := []string{"pl", "nl", "en"}
	mConfig, _ := config.NewAppConfig("", "", "", "", "", true)

	for _, lang := range langs {
		os.Setenv("LC_ALL", lang)
		mApp, _ := app.NewApp(mConfig, "")
		file, err := os.Create(getProjectRoot() + "/docs/keybindings/Keybindings_" + lang + ".md")
		if err != nil {
			panic(err)
		}

		bindingSections := getBindingSections(mApp)
		content := formatSections(mApp, bindingSections)
		writeString(file, content)
	}
}

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

func formatBinding(binding *gui.Binding) string {
	if binding.Alternative != "" {
		return fmt.Sprintf("  <kbd>%s</kbd>: %s (%s)\n", gui.GetKeyDisplay(binding.Key), binding.Description, binding.Alternative)
	}
	return fmt.Sprintf("  <kbd>%s</kbd>: %s\n", gui.GetKeyDisplay(binding.Key), binding.Description)
}

func getBindingSections(mApp *app.App) []*bindingSection {
	bindingSections := []*bindingSection{}

	bindings := mApp.Gui.GetInitialKeybindings()

	type contextAndViewType struct {
		context  string
		viewName string
	}

	contextAndViewBindingMap := map[contextAndViewType][]*gui.Binding{}

	for _, binding := range bindings {
		contexts := []string{}
		if len(binding.Contexts) == 0 {
			contexts = append(contexts, "")
		} else {
			contexts = append(contexts, binding.Contexts...)
		}

		for _, context := range contexts {
			key := contextAndViewType{context: context, viewName: binding.ViewName}
			existing := contextAndViewBindingMap[key]
			if existing == nil {
				contextAndViewBindingMap[key] = []*gui.Binding{binding}
			} else {
				contextAndViewBindingMap[key] = append(contextAndViewBindingMap[key], binding)
			}
		}
	}

	type groupedBindingsType struct {
		contextAndView contextAndViewType
		bindings       []*gui.Binding
	}

	groupedBindings := make([]groupedBindingsType, len(contextAndViewBindingMap))

	for contextAndView, contextBindings := range contextAndViewBindingMap {
		groupedBindings = append(groupedBindings, groupedBindingsType{contextAndView: contextAndView, bindings: contextBindings})
	}

	sort.Slice(groupedBindings, func(i, j int) bool {
		first := groupedBindings[i].contextAndView
		second := groupedBindings[j].contextAndView
		if first.viewName == "" {
			return true
		}
		if second.viewName == "" {
			return false
		}
		return first.viewName < second.viewName || (first.viewName == second.viewName && first.context < second.context)
	})

	for _, group := range groupedBindings {
		contextAndView := group.contextAndView
		contextBindings := group.bindings
		mApp.Log.Warn("viewname: " + contextAndView.viewName + ", context: " + contextAndView.context)
		viewName := contextAndView.viewName
		if viewName == "" {
			viewName = "global"
		}
		translatedView := localisedTitle(mApp, viewName)
		var title string
		if contextAndView.context == "" {
			addendum := " " + mApp.Tr.SLocalize("Panel")
			if viewName == "global" {
				addendum = ""
			}
			title = fmt.Sprintf("%s%s", translatedView, addendum)
		} else {
			translatedContextName := localisedTitle(mApp, contextAndView.context)
			title = fmt.Sprintf("%s %s (%s)", translatedView, mApp.Tr.SLocalize("Panel"), translatedContextName)
		}

		for _, binding := range contextBindings {
			bindingSections = addBinding(title, bindingSections, binding)
		}
	}

	return bindingSections
}

func addBinding(title string, bindingSections []*bindingSection, binding *gui.Binding) []*bindingSection {
	if binding.Description == "" && binding.Alternative == "" {
		return bindingSections
	}

	for _, section := range bindingSections {
		if title == section.title {
			section.bindings = append(section.bindings, binding)
			return bindingSections
		}
	}

	section := &bindingSection{
		title:    title,
		bindings: []*gui.Binding{binding},
	}

	return append(bindingSections, section)
}

func formatSections(mApp *app.App, bindingSections []*bindingSection) string {
	content := fmt.Sprintf("# Lazygit %s\n", mApp.Tr.SLocalize("Keybindings"))

	for _, section := range bindingSections {
		content += formatTitle(section.title)
		content += "<pre>\n"
		for _, binding := range section.bindings {
			content += formatBinding(binding)
		}
		content += "</pre>\n"
	}

	return content
}

func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return strings.Split(dir, "lazygit")[0] + "lazygit"
}
