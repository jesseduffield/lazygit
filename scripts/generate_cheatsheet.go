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
	tr := mApp.Tr

	contextTitleMap := map[string]string{
		"global":         tr.GlobalTitle,
		"navigation":     tr.NavigationTitle,
		"branches":       tr.BranchesTitle,
		"localBranches":  tr.LocalBranchesTitle,
		"files":          tr.FilesTitle,
		"status":         tr.StatusTitle,
		"submodules":     tr.SubmodulesTitle,
		"subCommits":     tr.SubCommitsTitle,
		"remoteBranches": tr.RemoteBranchesTitle,
		"remotes":        tr.RemotesTitle,
		"reflogCommits":  tr.ReflogCommitsTitle,
		"tags":           tr.TagsTitle,
		"commitFiles":    tr.CommitFilesTitle,
		"commitMessage":  tr.CommitMessageTitle,
		"commits":        tr.CommitsTitle,
		"confirmation":   tr.ConfirmationTitle,
		"credentials":    tr.CredentialsTitle,
		"information":    tr.InformationTitle,
		"main":           tr.MainTitle,
		"patchBuilding":  tr.PatchBuildingTitle,
		"merging":        tr.MergingTitle,
		"normal":         tr.NormalTitle,
		"staging":        tr.StagingTitle,
		"menu":           tr.MenuTitle,
		"search":         tr.SearchTitle,
		"secondary":      tr.SecondaryTitle,
		"stash":          tr.StashTitle,
	}

	title, ok := contextTitleMap[str]
	if !ok {
		panic(fmt.Sprintf("title not found for %s", str))
	}

	return title
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
		subtitle string
		title    string
	}

	contextAndViewBindingMap := map[contextAndViewType][]*gui.Binding{}

outer:
	for _, binding := range bindings {
		if binding.Tag == "navigation" {
			key := contextAndViewType{subtitle: "", title: "navigation"}
			existing := contextAndViewBindingMap[key]
			if existing == nil {
				contextAndViewBindingMap[key] = []*gui.Binding{binding}
			} else {
				for _, navBinding := range contextAndViewBindingMap[key] {
					if navBinding.Description == binding.Description {
						continue outer
					}
				}
				contextAndViewBindingMap[key] = append(contextAndViewBindingMap[key], binding)
			}

			continue outer
		}

		contexts := []string{}
		if len(binding.Contexts) == 0 {
			contexts = append(contexts, "")
		} else {
			contexts = append(contexts, binding.Contexts...)
		}

		for _, context := range contexts {
			key := contextAndViewType{subtitle: context, title: binding.ViewName}
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
		if first.title == "" {
			return true
		}
		if second.title == "" {
			return false
		}
		if first.title == "navigation" {
			return true
		}
		if second.title == "navigation" {
			return false
		}
		return first.title < second.title || (first.title == second.title && first.subtitle < second.subtitle)
	})

	for _, group := range groupedBindings {
		contextAndView := group.contextAndView
		contextBindings := group.bindings
		mApp.Log.Info("viewname: " + contextAndView.title + ", context: " + contextAndView.subtitle)
		viewName := contextAndView.title
		if viewName == "" {
			viewName = "global"
		}
		translatedView := localisedTitle(mApp, viewName)
		var title string
		if contextAndView.subtitle == "" {
			addendum := " " + mApp.Tr.Panel
			if viewName == "global" || viewName == "navigation" {
				addendum = ""
			}
			title = fmt.Sprintf("%s%s", translatedView, addendum)
		} else {
			translatedContextName := localisedTitle(mApp, contextAndView.subtitle)
			title = fmt.Sprintf("%s %s (%s)", translatedView, mApp.Tr.Panel, translatedContextName)
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
	content := fmt.Sprintf("# Lazygit %s\n", mApp.Tr.Keybindings)

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
