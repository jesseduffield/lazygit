package gui

import (
	"log"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) menuListContext() IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:        "menu",
			Key:             "menu",
			Kind:            PERSISTENT_POPUP,
			OnGetOptionsMap: gui.getMenuOptions,
		},
		GetItemsLength:      func() int { return gui.Views.Menu.LinesHeight() },
		OnGetPanelState:     func() IListPanelState { return gui.State.Panels.Menu },
		OnClickSelectedItem: gui.onMenuPress,
		Gui:                 gui,

		// no GetDisplayStrings field because we do a custom render on menu creation
	}
}

func (gui *Gui) filesListContext() IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "files",
			WindowName: "files",
			Key:        FILES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:      func() int { return gui.State.FileTreeViewModel.GetItemsLength() },
		OnGetPanelState:     func() IListPanelState { return gui.State.Panels.Files },
		OnFocus:             OnFocusWrapper(gui.onFocusFile),
		OnRenderToMain:      OnFocusWrapper(gui.filesRenderToMain),
		OnClickSelectedItem: gui.handleFilePress,
		Gui:                 gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			lines := presentation.RenderFileTree(gui.State.FileTreeViewModel, gui.State.Modes.Diffing.Ref, gui.State.Submodules)
			mappedLines := make([][]string, len(lines))
			for i, line := range lines {
				mappedLines[i] = []string{line}
			}

			return mappedLines
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedFileNode()
			return item, item != nil
		},
	}
}

func (gui *Gui) branchesListContext() IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        LOCAL_BRANCHES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.Branches) },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.Branches },
		OnRenderToMain:  OnFocusWrapper(gui.branchesRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetBranchListDisplayStrings(gui.State.Branches, gui.State.ScreenMode != SCREEN_NORMAL, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedBranch()
			return item, item != nil
		},
	}
}

func (gui *Gui) remotesListContext() IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        REMOTES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:      func() int { return len(gui.State.Remotes) },
		OnGetPanelState:     func() IListPanelState { return gui.State.Panels.Remotes },
		OnRenderToMain:      OnFocusWrapper(gui.remotesRenderToMain),
		OnClickSelectedItem: gui.handleRemoteEnter,
		Gui:                 gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetRemoteListDisplayStrings(gui.State.Remotes, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedRemote()
			return item, item != nil
		},
	}
}

func (gui *Gui) remoteBranchesListContext() IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        REMOTE_BRANCHES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.RemoteBranches) },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.RemoteBranches },
		OnRenderToMain:  OnFocusWrapper(gui.remoteBranchesRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetRemoteBranchListDisplayStrings(gui.State.RemoteBranches, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedRemoteBranch()
			return item, item != nil
		},
	}
}

func (gui *Gui) tagsListContext() IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        TAGS_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.Tags) },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.Tags },
		OnRenderToMain:  OnFocusWrapper(gui.tagsRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetTagListDisplayStrings(gui.State.Tags, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedTag()
			return item, item != nil
		},
	}
}

func (gui *Gui) branchCommitsListContext() IListContext {
	parseEmoji := gui.UserConfig.Git.ParseEmoji
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "commits",
			WindowName: "commits",
			Key:        BRANCH_COMMITS_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:      func() int { return len(gui.State.Commits) },
		OnGetPanelState:     func() IListPanelState { return gui.State.Panels.Commits },
		OnFocus:             OnFocusWrapper(gui.onCommitFocus),
		OnRenderToMain:      OnFocusWrapper(gui.branchCommitsRenderToMain),
		OnClickSelectedItem: gui.handleViewCommitFiles,
		Gui:                 gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			selectedCommitSha := ""
			if gui.currentContext().GetKey() == BRANCH_COMMITS_CONTEXT_KEY {
				selectedCommit := gui.getSelectedLocalCommit()
				if selectedCommit != nil {
					selectedCommitSha = selectedCommit.Sha
				}
			}
			return presentation.GetCommitListDisplayStrings(
				gui.State.Commits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.cherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
				parseEmoji,
				selectedCommitSha,
				startIdx,
				length,
				gui.shouldShowGraph(),
				gui.State.BisectInfo,
			)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedLocalCommit()
			return item, item != nil
		},
		RenderSelection: true,
	}
}

func (gui *Gui) subCommitsListContext() IListContext {
	parseEmoji := gui.UserConfig.Git.ParseEmoji
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        SUB_COMMITS_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.SubCommits) },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.SubCommits },
		OnRenderToMain:  OnFocusWrapper(gui.subCommitsRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			selectedCommitSha := ""
			if gui.currentContext().GetKey() == SUB_COMMITS_CONTEXT_KEY {
				selectedCommit := gui.getSelectedSubCommit()
				if selectedCommit != nil {
					selectedCommitSha = selectedCommit.Sha
				}
			}
			return presentation.GetCommitListDisplayStrings(
				gui.State.SubCommits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.cherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
				parseEmoji,
				selectedCommitSha,
				startIdx,
				length,
				gui.shouldShowGraph(),
				git_commands.NewNullBisectInfo(),
			)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedSubCommit()
			return item, item != nil
		},
		RenderSelection: true,
	}
}

func (gui *Gui) shouldShowGraph() bool {
	if gui.State.Modes.Filtering.Active() {
		return false
	}

	value := gui.UserConfig.Git.Log.ShowGraph
	switch value {
	case "always":
		return true
	case "never":
		return false
	case "when-maximised":
		return gui.State.ScreenMode != SCREEN_NORMAL
	}

	log.Fatalf("Unknown value for git.log.showGraph: %s. Expected one of: 'always', 'never', 'when-maximised'", value)
	return false
}

func (gui *Gui) reflogCommitsListContext() IListContext {
	parseEmoji := gui.UserConfig.Git.ParseEmoji
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "commits",
			WindowName: "commits",
			Key:        REFLOG_COMMITS_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.FilteredReflogCommits) },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.ReflogCommits },
		OnRenderToMain:  OnFocusWrapper(gui.reflogCommitsRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetReflogCommitListDisplayStrings(
				gui.State.FilteredReflogCommits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.cherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
				parseEmoji,
			)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedReflogCommit()
			return item, item != nil
		},
	}
}

func (gui *Gui) stashListContext() IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "stash",
			WindowName: "stash",
			Key:        STASH_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.StashEntries) },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.Stash },
		OnRenderToMain:  OnFocusWrapper(gui.stashRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetStashEntryListDisplayStrings(gui.State.StashEntries, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedStashEntry()
			return item, item != nil
		},
	}
}

func (gui *Gui) commitFilesListContext() IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "commitFiles",
			WindowName: "commits",
			Key:        COMMIT_FILES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return gui.State.CommitFileTreeViewModel.GetItemsLength() },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.CommitFiles },
		OnFocus:         OnFocusWrapper(gui.onCommitFileFocus),
		OnRenderToMain:  OnFocusWrapper(gui.commitFilesRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			if gui.State.CommitFileTreeViewModel.GetItemsLength() == 0 {
				return [][]string{{style.FgRed.Sprint("(none)")}}
			}

			lines := presentation.RenderCommitFileTree(gui.State.CommitFileTreeViewModel, gui.State.Modes.Diffing.Ref, gui.Git.Patch.PatchManager)
			mappedLines := make([][]string, len(lines))
			for i, line := range lines {
				mappedLines[i] = []string{line}
			}

			return mappedLines
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedCommitFileNode()
			return item, item != nil
		},
	}
}

func (gui *Gui) submodulesListContext() IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "files",
			WindowName: "files",
			Key:        SUBMODULES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.Submodules) },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.Submodules },
		OnRenderToMain:  OnFocusWrapper(gui.submodulesRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetSubmoduleListDisplayStrings(gui.State.Submodules)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedSubmodule()
			return item, item != nil
		},
	}
}

func (gui *Gui) suggestionsListContext() IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "suggestions",
			WindowName: "suggestions",
			Key:        SUGGESTIONS_CONTEXT_KEY,
			Kind:       PERSISTENT_POPUP,
		},
		GetItemsLength:  func() int { return len(gui.State.Suggestions) },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.Suggestions },
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetSuggestionListDisplayStrings(gui.State.Suggestions)
		},
	}
}

func (gui *Gui) getListContexts() []IListContext {
	return []IListContext{
		gui.State.Contexts.Menu,
		gui.State.Contexts.Files,
		gui.State.Contexts.Branches,
		gui.State.Contexts.Remotes,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Tags,
		gui.State.Contexts.BranchCommits,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.SubCommits,
		gui.State.Contexts.Stash,
		gui.State.Contexts.CommitFiles,
		gui.State.Contexts.Submodules,
		gui.State.Contexts.Suggestions,
	}
}

func (gui *Gui) getListContextKeyBindings() []*types.Binding {
	bindings := make([]*types.Binding, 0)

	keybindingConfig := gui.UserConfig.Keybinding

	for _, listContext := range gui.getListContexts() {
		listContext := listContext

		bindings = append(bindings, []*types.Binding{
			{ViewName: listContext.GetViewName(), Tag: "navigation", Contexts: []string{string(listContext.GetKey())}, Key: gui.getKey(keybindingConfig.Universal.PrevItemAlt), Modifier: gocui.ModNone, Handler: listContext.handlePrevLine},
			{ViewName: listContext.GetViewName(), Tag: "navigation", Contexts: []string{string(listContext.GetKey())}, Key: gui.getKey(keybindingConfig.Universal.PrevItem), Modifier: gocui.ModNone, Handler: listContext.handlePrevLine},
			{ViewName: listContext.GetViewName(), Tag: "navigation", Contexts: []string{string(listContext.GetKey())}, Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: listContext.handlePrevLine},
			{ViewName: listContext.GetViewName(), Tag: "navigation", Contexts: []string{string(listContext.GetKey())}, Key: gui.getKey(keybindingConfig.Universal.NextItemAlt), Modifier: gocui.ModNone, Handler: listContext.handleNextLine},
			{ViewName: listContext.GetViewName(), Tag: "navigation", Contexts: []string{string(listContext.GetKey())}, Key: gui.getKey(keybindingConfig.Universal.NextItem), Modifier: gocui.ModNone, Handler: listContext.handleNextLine},
			{ViewName: listContext.GetViewName(), Tag: "navigation", Contexts: []string{string(listContext.GetKey())}, Key: gui.getKey(keybindingConfig.Universal.PrevPage), Modifier: gocui.ModNone, Handler: listContext.handlePrevPage, Description: gui.Tr.LcPrevPage},
			{ViewName: listContext.GetViewName(), Tag: "navigation", Contexts: []string{string(listContext.GetKey())}, Key: gui.getKey(keybindingConfig.Universal.NextPage), Modifier: gocui.ModNone, Handler: listContext.handleNextPage, Description: gui.Tr.LcNextPage},
			{ViewName: listContext.GetViewName(), Tag: "navigation", Contexts: []string{string(listContext.GetKey())}, Key: gui.getKey(keybindingConfig.Universal.GotoTop), Modifier: gocui.ModNone, Handler: listContext.handleGotoTop, Description: gui.Tr.LcGotoTop},
			{ViewName: listContext.GetViewName(), Tag: "navigation", Contexts: []string{string(listContext.GetKey())}, Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: listContext.handleNextLine},
			{ViewName: listContext.GetViewName(), Contexts: []string{string(listContext.GetKey())}, Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: listContext.handleClick},
			{ViewName: listContext.GetViewName(), Tag: "navigation", Contexts: []string{string(listContext.GetKey())}, Key: gui.getKey(keybindingConfig.Universal.ScrollLeft), Modifier: gocui.ModNone, Handler: listContext.handleScrollLeft},
			{ViewName: listContext.GetViewName(), Tag: "navigation", Contexts: []string{string(listContext.GetKey())}, Key: gui.getKey(keybindingConfig.Universal.ScrollRight), Modifier: gocui.ModNone, Handler: listContext.handleScrollRight},
		}...)

		openSearchHandler := gui.handleOpenSearch
		gotoBottomHandler := listContext.handleGotoBottom

		// the branch commits context needs to lazyload things so it has a couple of its own handlers
		if listContext.GetKey() == BRANCH_COMMITS_CONTEXT_KEY {
			openSearchHandler = gui.handleOpenSearchForCommitsPanel
			gotoBottomHandler = gui.handleGotoBottomForCommitsPanel
		}

		bindings = append(bindings, []*types.Binding{
			{
				ViewName:    listContext.GetViewName(),
				Contexts:    []string{string(listContext.GetKey())},
				Key:         gui.getKey(keybindingConfig.Universal.StartSearch),
				Handler:     func() error { return openSearchHandler(listContext.GetViewName()) },
				Description: gui.Tr.LcStartSearch,
				Tag:         "navigation",
			},
			{
				ViewName:    listContext.GetViewName(),
				Contexts:    []string{string(listContext.GetKey())},
				Key:         gui.getKey(keybindingConfig.Universal.GotoBottom),
				Handler:     gotoBottomHandler,
				Description: gui.Tr.LcGotoBottom,
				Tag:         "navigation",
			},
		}...)
	}

	return bindings
}
