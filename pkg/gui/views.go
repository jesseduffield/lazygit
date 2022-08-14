package gui

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

type Views struct {
	Status         *gocui.View
	Submodules     *gocui.View
	Files          *gocui.View
	Branches       *gocui.View
	Remotes        *gocui.View
	Tags           *gocui.View
	RemoteBranches *gocui.View
	ReflogCommits  *gocui.View
	Commits        *gocui.View
	Stash          *gocui.View

	Main                   *gocui.View
	Secondary              *gocui.View
	Staging                *gocui.View
	StagingSecondary       *gocui.View
	PatchBuilding          *gocui.View
	PatchBuildingSecondary *gocui.View
	MergeConflicts         *gocui.View

	Options       *gocui.View
	Confirmation  *gocui.View
	Menu          *gocui.View
	CommitMessage *gocui.View
	CommitFiles   *gocui.View
	SubCommits    *gocui.View
	Information   *gocui.View
	AppStatus     *gocui.View
	Search        *gocui.View
	SearchPrefix  *gocui.View
	Limit         *gocui.View
	Suggestions   *gocui.View
	Tooltip       *gocui.View
	Extras        *gocui.View
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
		{viewPtr: &gui.Views.Submodules, name: "submodules"},
		{viewPtr: &gui.Views.Files, name: "files"},
		{viewPtr: &gui.Views.Tags, name: "tags"},
		{viewPtr: &gui.Views.Remotes, name: "remotes"},
		{viewPtr: &gui.Views.Branches, name: "localBranches"},
		{viewPtr: &gui.Views.RemoteBranches, name: "remoteBranches"},
		{viewPtr: &gui.Views.ReflogCommits, name: "reflogCommits"},
		{viewPtr: &gui.Views.Commits, name: "commits"},
		{viewPtr: &gui.Views.Stash, name: "stash"},
		{viewPtr: &gui.Views.SubCommits, name: "subCommits"},
		{viewPtr: &gui.Views.CommitFiles, name: "commitFiles"},

		{viewPtr: &gui.Views.Staging, name: "staging"},
		{viewPtr: &gui.Views.StagingSecondary, name: "stagingSecondary"},
		{viewPtr: &gui.Views.PatchBuilding, name: "patchBuilding"},
		{viewPtr: &gui.Views.PatchBuildingSecondary, name: "patchBuildingSecondary"},
		{viewPtr: &gui.Views.MergeConflicts, name: "mergeConflicts"},
		{viewPtr: &gui.Views.Secondary, name: "secondary"},
		{viewPtr: &gui.Views.Main, name: "main"},

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

func (gui *Gui) windowForView(viewName string) string {
	context, ok := gui.contextForView(viewName)
	if !ok {
		panic("todo: deal with this")
	}

	return context.GetWindowName()
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
	gui.Views.Options.Frame = false

	gui.Views.SearchPrefix.BgColor = gocui.ColorDefault
	gui.Views.SearchPrefix.FgColor = gocui.ColorGreen
	gui.Views.SearchPrefix.Frame = false
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

	gui.Views.Remotes.Title = gui.c.Tr.RemotesTitle
	gui.Views.Remotes.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Tags.Title = gui.c.Tr.TagsTitle
	gui.Views.Tags.FgColor = theme.GocuiDefaultTextColor

	gui.Views.RemoteBranches.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Files.Title = gui.c.Tr.FilesTitle
	gui.Views.Files.FgColor = theme.GocuiDefaultTextColor

	for _, view := range []*gocui.View{gui.Views.Main, gui.Views.Secondary, gui.Views.Staging, gui.Views.StagingSecondary, gui.Views.PatchBuilding, gui.Views.PatchBuildingSecondary, gui.Views.MergeConflicts} {
		view.Title = gui.c.Tr.DiffTitle
		view.Wrap = true
		view.FgColor = theme.GocuiDefaultTextColor
		view.IgnoreCarriageReturns = true
		view.CanScrollPastBottom = gui.c.UserConfig.Gui.ScrollPastBottom
	}

	gui.Views.Staging.Title = gui.c.Tr.UnstagedChanges
	gui.Views.Staging.Highlight = true
	gui.Views.Staging.Wrap = true

	gui.Views.StagingSecondary.Title = gui.c.Tr.StagedChanges
	gui.Views.StagingSecondary.Highlight = true
	gui.Views.StagingSecondary.Wrap = true

	gui.Views.PatchBuilding.Title = gui.Tr.Patch
	gui.Views.PatchBuilding.Highlight = true
	gui.Views.PatchBuilding.Wrap = true

	gui.Views.PatchBuildingSecondary.Title = gui.Tr.CustomPatch
	gui.Views.PatchBuildingSecondary.Highlight = true
	gui.Views.PatchBuildingSecondary.Wrap = true

	gui.Views.MergeConflicts.Title = gui.c.Tr.MergeConflictsTitle
	gui.Views.MergeConflicts.Highlight = true
	gui.Views.MergeConflicts.Wrap = false

	gui.Views.Limit.Title = gui.c.Tr.NotEnoughSpace
	gui.Views.Limit.Wrap = true

	gui.Views.Status.Title = gui.c.Tr.StatusTitle
	gui.Views.Status.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Search.BgColor = gocui.ColorDefault
	gui.Views.Search.FgColor = gocui.ColorGreen
	gui.Views.Search.Editable = true
	gui.Views.Search.Frame = false

	gui.Views.AppStatus.BgColor = gocui.ColorDefault
	gui.Views.AppStatus.FgColor = gocui.ColorCyan
	gui.Views.AppStatus.Visible = false
	gui.Views.AppStatus.Frame = false

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
	gui.Views.Information.Frame = false

	gui.Views.Extras.Title = gui.c.Tr.CommandLog
	gui.Views.Extras.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Extras.Autoscroll = true
	gui.Views.Extras.Wrap = true

	return nil
}
