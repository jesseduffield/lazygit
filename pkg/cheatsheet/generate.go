//go:generate go run generator.go

// This "script" generates files called Keybindings_{{.LANG}}.md
// in the docs/keybindings directory.
//
// The content of these generated files is a keybindings cheatsheet.
//
// To generate the cheatsheets, run:
//   go generate pkg/cheatsheet/generate.go

package cheatsheet

import (
	"fmt"
	"log"
	"os"

	"github.com/jesseduffield/generics/maps"
	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type bindingSection struct {
	title    string
	bindings []*types.Binding
}

type header struct {
	// priority decides the order of the headers in the cheatsheet (lower means higher)
	priority int
	title    string
}

type headerWithBindings struct {
	header   header
	bindings []*types.Binding
}

func CommandToRun() string {
	return "go generate ./..."
}

func GetKeybindingsDir() string {
	return utils.GetLazyRootDirectory() + "/docs/keybindings"
}

func generateAtDir(cheatsheetDir string) {
	translationSetsByLang, err := i18n.GetTranslationSets()
	if err != nil {
		log.Fatal(err)
	}
	mConfig := config.NewDummyAppConfig()

	for lang := range translationSetsByLang {
		mConfig.GetUserConfig().Gui.Language = lang
		common, err := app.NewCommon(mConfig)
		if err != nil {
			log.Fatal(err)
		}
		tr, err := i18n.NewTranslationSetFromConfig(common.Log, lang)
		if err != nil {
			log.Fatal(err)
		}
		common.Tr = tr
		mApp, _ := app.NewApp(mConfig, nil, common)
		path := cheatsheetDir + "/Keybindings_" + lang + ".md"
		file, err := os.Create(path)
		if err != nil {
			panic(err)
		}

		bindings := mApp.Gui.GetCheatsheetKeybindings()
		bindingSections := getBindingSections(bindings, mApp.Tr)
		content := formatSections(mApp.Tr, bindingSections)
		content = fmt.Sprintf("_This file is auto-generated. To update, make the changes in the "+
			"pkg/i18n directory and then run `%s` from the project root._\n\n%s", CommandToRun(), content)
		writeString(file, content)
	}
}

func Generate() {
	generateAtDir(GetKeybindingsDir())
}

func writeString(file *os.File, str string) {
	_, err := file.WriteString(str)
	if err != nil {
		log.Fatal(err)
	}
}

func localisedTitle(tr *i18n.TranslationSet, str string) string {
	contextTitleMap := map[string]string{
		"global":            tr.GlobalTitle,
		"navigation":        tr.NavigationTitle,
		"branches":          tr.BranchesTitle,
		"localBranches":     tr.LocalBranchesTitle,
		"files":             tr.FilesTitle,
		"status":            tr.StatusTitle,
		"submodules":        tr.SubmodulesTitle,
		"subCommits":        tr.SubCommitsTitle,
		"remoteBranches":    tr.RemoteBranchesTitle,
		"remotes":           tr.RemotesTitle,
		"reflogCommits":     tr.ReflogCommitsTitle,
		"tags":              tr.TagsTitle,
		"commitFiles":       tr.CommitFilesTitle,
		"commitMessage":     tr.CommitSummaryTitle,
		"commitDescription": tr.CommitDescriptionTitle,
		"commits":           tr.CommitsTitle,
		"confirmation":      tr.ConfirmationTitle,
		"information":       tr.InformationTitle,
		"main":              tr.NormalTitle,
		"patchBuilding":     tr.PatchBuildingTitle,
		"mergeConflicts":    tr.MergingTitle,
		"staging":           tr.StagingTitle,
		"menu":              tr.MenuTitle,
		"search":            tr.SearchTitle,
		"secondary":         tr.SecondaryTitle,
		"stash":             tr.StashTitle,
		"suggestions":       tr.SuggestionsCheatsheetTitle,
		"extras":            tr.ExtrasTitle,
		"worktrees":         tr.WorktreesTitle,
	}

	title, ok := contextTitleMap[str]
	if !ok {
		panic(fmt.Sprintf("title not found for %s", str))
	}

	return title
}

func getBindingSections(bindings []*types.Binding, tr *i18n.TranslationSet) []*bindingSection {
	excludedViews := []string{"stagingSecondary", "patchBuildingSecondary"}
	bindingsToDisplay := lo.Filter(bindings, func(binding *types.Binding, _ int) bool {
		if lo.Contains(excludedViews, binding.ViewName) {
			return false
		}

		return (binding.Description != "" || binding.Alternative != "") && binding.Key != nil
	})

	bindingsByHeader := lo.GroupBy(bindingsToDisplay, func(binding *types.Binding) header {
		return getHeader(binding, tr)
	})

	bindingGroups := maps.MapToSlice(
		bindingsByHeader,
		func(header header, hBindings []*types.Binding) headerWithBindings {
			uniqBindings := lo.UniqBy(hBindings, func(binding *types.Binding) string {
				return binding.Description + keybindings.LabelFromKey(binding.Key)
			})

			return headerWithBindings{
				header:   header,
				bindings: uniqBindings,
			}
		},
	)

	slices.SortFunc(bindingGroups, func(a, b headerWithBindings) bool {
		if a.header.priority != b.header.priority {
			return a.header.priority > b.header.priority
		}
		return a.header.title < b.header.title
	})

	return lo.Map(bindingGroups, func(hb headerWithBindings, _ int) *bindingSection {
		return &bindingSection{
			title:    hb.header.title,
			bindings: hb.bindings,
		}
	})
}

func getHeader(binding *types.Binding, tr *i18n.TranslationSet) header {
	if binding.Tag == "navigation" {
		return header{priority: 2, title: localisedTitle(tr, "navigation")}
	}

	if binding.ViewName == "" {
		return header{priority: 3, title: localisedTitle(tr, "global")}
	}

	return header{priority: 1, title: localisedTitle(tr, binding.ViewName)}
}

func formatSections(tr *i18n.TranslationSet, bindingSections []*bindingSection) string {
	content := fmt.Sprintf("# Lazygit %s\n", tr.Keybindings)

	content += fmt.Sprintf("\n%s\n", italicize(tr.KeybindingsLegend))

	for _, section := range bindingSections {
		content += formatTitle(section.title)
		content += "| Key | Action | Info |\n"
		content += "|-----|--------|-------------|\n"
		for _, binding := range section.bindings {
			content += formatBinding(binding)
		}
	}

	return content
}

func formatTitle(title string) string {
	return fmt.Sprintf("\n## %s\n\n", title)
}

func formatBinding(binding *types.Binding) string {
	action := keybindings.LabelFromKey(binding.Key)
	description := binding.Description
	if binding.Alternative != "" {
		action += fmt.Sprintf(" (%s)", binding.Alternative)
	}

	// Use backticks for keyboard keys. Two backticks are needed with an inner space
	//  to escape a key that is itself a backtick.
	return fmt.Sprintf("| `` %s `` | %s | %s |\n", action, description, binding.Tooltip)
}

func italicize(str string) string {
	return fmt.Sprintf("_%s_", str)
}
