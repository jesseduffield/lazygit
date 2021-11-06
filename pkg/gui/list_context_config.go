package gui

import (
	"log"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
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
		OnFocus:             gui.handleMenuSelect,
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
		GetItemsLength:      func() int { return gui.State.FileManager.GetItemsLength() },
		OnGetPanelState:     func() IListPanelState { return gui.State.Panels.Files },
		OnFocus:             gui.focusAndSelectFile,
		OnClickSelectedItem: gui.handleFilePress,
		Gui:                 gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			lines := gui.State.FileManager.Render(gui.State.Modes.Diffing.Ref, gui.State.Submodules)
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
		OnFocus:         gui.handleBranchSelect,
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			prs := gui.GitCommand.GenerateGithubPullRequestMap(gui.State.GithubState.RecentPRs, gui.State.Branches)
			return presentation.GetBranchListDisplayStrings(gui.State.Branches, prs, gui.State.ScreenMode != SCREEN_NORMAL, gui.State.Modes.Diffing.Ref)
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
		OnFocus:             gui.handleRemoteSelect,
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
		OnFocus:         gui.handleRemoteBranchSelect,
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
		OnFocus:         gui.handleTagSelect,
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
	parseEmoji := gui.Config.GetUserConfig().Git.ParseEmoji
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "commits",
			WindowName: "commits",
			Key:        BRANCH_COMMITS_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:      func() int { return len(gui.State.Commits) },
		OnGetPanelState:     func() IListPanelState { return gui.State.Panels.Commits },
		OnFocus:             gui.handleCommitSelect,
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
			)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedLocalCommit()
			return item, item != nil
		},
		RenderSelection: true,
	}
}

func (gui *Gui) shouldShowGraph() bool {
	value := gui.Config.GetUserConfig().Git.Log.ShowGraph
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
	parseEmoji := gui.Config.GetUserConfig().Git.ParseEmoji
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "commits",
			WindowName: "commits",
			Key:        REFLOG_COMMITS_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.FilteredReflogCommits) },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.ReflogCommits },
		OnFocus:         gui.handleReflogCommitSelect,
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

func (gui *Gui) subCommitsListContext() IListContext {
	parseEmoji := gui.Config.GetUserConfig().Git.ParseEmoji
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        SUB_COMMITS_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.SubCommits) },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.SubCommits },
		OnFocus:         gui.handleSubCommitSelect,
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
				0,
				len(gui.State.SubCommits),
				gui.shouldShowGraph(),
			)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedSubCommit()
			return item, item != nil
		},
		RenderSelection: true,
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
		OnFocus:         gui.handleStashEntrySelect,
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
		GetItemsLength:  func() int { return gui.State.CommitFileManager.GetItemsLength() },
		OnGetPanelState: func() IListPanelState { return gui.State.Panels.CommitFiles },
		OnFocus:         gui.handleCommitFileSelect,
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			if gui.State.CommitFileManager.GetItemsLength() == 0 {
				return [][]string{{style.FgRed.Sprint("(none)")}}
			}

			lines := gui.State.CommitFileManager.Render(gui.State.Modes.Diffing.Ref, gui.GitCommand.PatchManager)
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
		OnFocus:         gui.handleSubmoduleSelect,
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
		OnFocus:         func() error { return nil },
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

func (gui *Gui) getListContextKeyBindings() []*Binding {
	bindings := make([]*Binding, 0)

	keybindingConfig := gui.Config.GetUserConfig().Keybinding

	for _, listContext := range gui.getListContexts() {
		listContext := listContext

		bindings = append(bindings, []*Binding{
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

		// the commits panel needs to lazyload things so it has a couple of its own handlers
		openSearchHandler := gui.handleOpenSearch
		gotoBottomHandler := listContext.handleGotoBottom
		if listContext.GetViewName() == "commits" {
			openSearchHandler = gui.handleOpenSearchForCommitsPanel
			gotoBottomHandler = gui.handleGotoBottomForCommitsPanel
		}

		bindings = append(bindings, []*Binding{
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
