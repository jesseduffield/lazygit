package gui

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

type Views struct {
	Status         *gocui.View
	Files          *gocui.View
	Branches       *gocui.View
	RemoteBranches *gocui.View
	Commits        *gocui.View
	Stash          *gocui.View
	Main           *gocui.View
	Secondary      *gocui.View
	Options        *gocui.View
	Confirmation   *gocui.View
	Menu           *gocui.View
	CommitMessage  *gocui.View
	CommitFiles    *gocui.View
	SubCommits     *gocui.View
	Information    *gocui.View
	AppStatus      *gocui.View
	Search         *gocui.View
	SearchPrefix   *gocui.View
	Limit          *gocui.View
	Suggestions    *gocui.View
	Tooltip        *gocui.View
	Extras         *gocui.View
}

type viewNameMapping struct {
	viewPtr **gocui.View
	name    string
}

func (gui *Gui) orderedViews() []*gocui.View {
	return slices.Map(gui.orderedViewNameMappings(), func(v viewNameMapping) *gocui.View {
		return *v.viewPtr
	})
}

func (gui *Gui) orderedViewNameMappings() []viewNameMapping {
	return []viewNameMapping{
		// first layer. Ordering within this layer does not matter because there are
		// no overlapping views
		{viewPtr: &gui.Views.Status, name: "status"},
		{viewPtr: &gui.Views.Files, name: "files"},
		{viewPtr: &gui.Views.Branches, name: "branches"},
		{viewPtr: &gui.Views.RemoteBranches, name: "remoteBranches"},
		{viewPtr: &gui.Views.Commits, name: "commits"},
		{viewPtr: &gui.Views.Stash, name: "stash"},
		{viewPtr: &gui.Views.SubCommits, name: "subCommits"},
		{viewPtr: &gui.Views.CommitFiles, name: "commitFiles"},
		{viewPtr: &gui.Views.Main, name: "main"},
		{viewPtr: &gui.Views.Secondary, name: "secondary"},
		{viewPtr: &gui.Views.Extras, name: "extras"},

		// bottom line
		{viewPtr: &gui.Views.Options, name: "options"},
		{viewPtr: &gui.Views.AppStatus, name: "appStatus"},
		{viewPtr: &gui.Views.Information, name: "information"},
		{viewPtr: &gui.Views.Search, name: "search"},
		// this view takes up one character. Its only purpose is to show the slash when searching
		{viewPtr: &gui.Views.SearchPrefix, name: "searchPrefix"},

		// popups.
		{viewPtr: &gui.Views.CommitMessage, name: "commitMessage"},
		{viewPtr: &gui.Views.Menu, name: "menu"},
		{viewPtr: &gui.Views.Suggestions, name: "suggestions"},
		{viewPtr: &gui.Views.Confirmation, name: "confirmation"},
		{viewPtr: &gui.Views.Tooltip, name: "tooltip"},

		// this guy will cover everything else when it appears
		{viewPtr: &gui.Views.Limit, name: "limit"},
	}
}

type controlledView struct {
	viewName   string
	windowName string
	frame      bool
}

// controlled views have their size and position determined in arrangement.go.
// Some views, like the confirmation panel, are currently sized at the time of
// displaying the view, based on the view's contents.
func (gui *Gui) controlledViews() []controlledView {
	return []controlledView{
		{viewName: "main", windowName: "main", frame: true},
		{viewName: "secondary", windowName: "secondary", frame: true},
		{viewName: "status", windowName: "status", frame: true},
		{viewName: "files", windowName: "files", frame: true},
		{viewName: "branches", windowName: "branches", frame: true},
		{viewName: "remoteBranches", windowName: "branches", frame: true},
		{viewName: "commitFiles", windowName: gui.State.Contexts.CommitFiles.GetWindowName(), frame: true},
		{viewName: "subCommits", windowName: gui.State.Contexts.SubCommits.GetWindowName(), frame: true},
		{viewName: "commits", windowName: "commits", frame: true},
		{viewName: "stash", windowName: "stash", frame: true},
		{viewName: "options", windowName: "options", frame: false},
		{viewName: "searchPrefix", windowName: "searchPrefix", frame: false},
		{viewName: "search", windowName: "search", frame: false},
		{viewName: "appStatus", windowName: "appStatus", frame: false},
		{viewName: "information", windowName: "information", frame: false},
		{viewName: "extras", windowName: "extras", frame: true},
		{viewName: "limit", windowName: "limit", frame: true},
	}
}

func (gui *Gui) createAllViews() error {
	var err error
	for _, mapping := range gui.orderedViewNameMappings() {
		*mapping.viewPtr, err = gui.prepareView(mapping.name)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
	}

	gui.Views.Options.FgColor = theme.OptionsColor

	gui.Views.SearchPrefix.BgColor = gocui.ColorDefault
	gui.Views.SearchPrefix.FgColor = gocui.ColorGreen
	gui.setViewContent(gui.Views.SearchPrefix, SEARCH_PREFIX)

	gui.Views.Stash.Title = gui.c.Tr.StashTitle
	gui.Views.Stash.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Commits.Title = gui.c.Tr.CommitsTitle
	gui.Views.Commits.FgColor = theme.GocuiDefaultTextColor

	gui.Views.CommitFiles.Title = gui.c.Tr.CommitFiles
	gui.Views.CommitFiles.FgColor = theme.GocuiDefaultTextColor

	gui.Views.SubCommits.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Branches.Title = gui.c.Tr.BranchesTitle
	gui.Views.Branches.FgColor = theme.GocuiDefaultTextColor

	gui.Views.RemoteBranches.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Files.Title = gui.c.Tr.FilesTitle
	gui.Views.Files.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Secondary.Title = gui.c.Tr.DiffTitle
	gui.Views.Secondary.Wrap = true
	gui.Views.Secondary.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Secondary.IgnoreCarriageReturns = true
	gui.Views.Secondary.CanScrollPastBottom = gui.c.UserConfig.Gui.ScrollPastBottom

	gui.Views.Main.Title = gui.c.Tr.DiffTitle
	gui.Views.Main.Wrap = true
	gui.Views.Main.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Main.IgnoreCarriageReturns = true
	gui.Views.Main.CanScrollPastBottom = gui.c.UserConfig.Gui.ScrollPastBottom

	gui.Views.Limit.Title = gui.c.Tr.NotEnoughSpace
	gui.Views.Limit.Wrap = true

	gui.Views.Status.Title = gui.c.Tr.StatusTitle
	gui.Views.Status.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Search.BgColor = gocui.ColorDefault
	gui.Views.Search.FgColor = gocui.ColorGreen
	gui.Views.Search.Editable = true

	gui.Views.AppStatus.BgColor = gocui.ColorDefault
	gui.Views.AppStatus.FgColor = gocui.ColorCyan
	gui.Views.AppStatus.Visible = false

	gui.Views.CommitMessage.Visible = false
	gui.Views.CommitMessage.Title = gui.c.Tr.CommitMessage
	gui.Views.CommitMessage.FgColor = theme.GocuiDefaultTextColor
	gui.Views.CommitMessage.Editable = true
	gui.Views.CommitMessage.Editor = gocui.EditorFunc(gui.commitMessageEditor)

	gui.Views.Confirmation.Visible = false

	gui.Views.Suggestions.Visible = false

	gui.Views.Tooltip.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Menu.Visible = false

	gui.Views.Tooltip.Visible = false

	gui.Views.Information.BgColor = gocui.ColorDefault
	gui.Views.Information.FgColor = gocui.ColorGreen

	gui.Views.Extras.Title = gui.c.Tr.CommandLog
	gui.Views.Extras.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Extras.Autoscroll = true
	gui.Views.Extras.Wrap = true

	return nil
}

func initialViewContextMapping(contextTree *context.ContextTree) map[string]types.Context {
	return map[string]types.Context{
		"status":         contextTree.Status,
		"files":          contextTree.Files,
		"branches":       contextTree.Branches,
		"remoteBranches": contextTree.RemoteBranches,
		"commits":        contextTree.LocalCommits,
		"commitFiles":    contextTree.CommitFiles,
		"subCommits":     contextTree.SubCommits,
		"stash":          contextTree.Stash,
		"menu":           contextTree.Menu,
		"confirmation":   contextTree.Confirmation,
		"commitMessage":  contextTree.CommitMessage,
		"main":           contextTree.Normal,
		"secondary":      contextTree.Normal,
		"extras":         contextTree.CommandLog,
	}
}
