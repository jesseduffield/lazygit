package gui

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) menuListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:        "menu",
			Key:             "menu",
			Kind:            PERSISTENT_POPUP,
			OnGetOptionsMap: gui.getMenuOptions,
		},
		GetItemsLength:             func() int { return gui.Views.Menu.LinesHeight() },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Menu },
		OnFocus:                    gui.handleMenuSelect,
		OnClickSelectedItem:        gui.onMenuPress,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: false,

		// no GetDisplayStrings field because we do a custom render on menu creation
	}
}

func (gui *Gui) filesListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "files",
			WindowName: "files",
			Key:        FILES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:             func() int { return gui.State.FileManager.GetItemsLength() },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Files },
		OnFocus:                    gui.focusAndSelectFile,
		OnClickSelectedItem:        gui.handleFilePress,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: false,
		GetDisplayStrings: func() [][]string {
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

func (gui *Gui) branchesListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        LOCAL_BRANCHES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:             func() int { return len(gui.State.Branches) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Branches },
		OnFocus:                    gui.handleBranchSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		GetDisplayStrings: func() [][]string {
			return presentation.GetBranchListDisplayStrings(gui.State.Branches, gui.State.ScreenMode != SCREEN_NORMAL, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedBranch()
			return item, item != nil
		},
	}
}

func (gui *Gui) remotesListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        REMOTES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:             func() int { return len(gui.State.Remotes) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Remotes },
		OnFocus:                    gui.handleRemoteSelect,
		OnClickSelectedItem:        gui.handleRemoteEnter,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		GetDisplayStrings: func() [][]string {
			return presentation.GetRemoteListDisplayStrings(gui.State.Remotes, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedRemote()
			return item, item != nil
		},
	}
}

func (gui *Gui) remoteBranchesListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        REMOTE_BRANCHES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:             func() int { return len(gui.State.RemoteBranches) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.RemoteBranches },
		OnFocus:                    gui.handleRemoteBranchSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		GetDisplayStrings: func() [][]string {
			return presentation.GetRemoteBranchListDisplayStrings(gui.State.RemoteBranches, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedRemoteBranch()
			return item, item != nil
		},
	}
}

func (gui *Gui) tagsListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        TAGS_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:             func() int { return len(gui.State.Tags) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Tags },
		OnFocus:                    gui.handleTagSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		GetDisplayStrings: func() [][]string {
			return presentation.GetTagListDisplayStrings(gui.State.Tags, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedTag()
			return item, item != nil
		},
	}
}

func (gui *Gui) branchCommitsListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "commits",
			WindowName: "commits",
			Key:        BRANCH_COMMITS_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:             func() int { return len(gui.State.Commits) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Commits },
		OnFocus:                    gui.handleCommitSelect,
		OnClickSelectedItem:        gui.handleViewCommitFiles,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		GetDisplayStrings: func() [][]string {
			return presentation.GetCommitListDisplayStrings(
				gui.State.Commits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.cherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
				gui.Git.Status().IsRebasing(),
				gui.Tr,
			)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedLocalCommit()
			return item, item != nil
		},
	}
}

func (gui *Gui) reflogCommitsListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "commits",
			WindowName: "commits",
			Key:        REFLOG_COMMITS_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:             func() int { return len(gui.State.FilteredReflogCommits) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.ReflogCommits },
		OnFocus:                    gui.handleReflogCommitSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		GetDisplayStrings: func() [][]string {
			return presentation.GetReflogCommitListDisplayStrings(
				gui.State.FilteredReflogCommits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.cherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
			)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedReflogCommit()
			return item, item != nil
		},
	}
}

func (gui *Gui) subCommitsListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        SUB_COMMITS_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:             func() int { return len(gui.State.SubCommits) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.SubCommits },
		OnFocus:                    gui.handleSubCommitSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		GetDisplayStrings: func() [][]string {
			return presentation.GetCommitListDisplayStrings(
				gui.State.SubCommits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.cherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
				// no need to mention if we're rebasing because that only applies to main commits panel
				false,
				gui.Tr,
			)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedSubCommit()
			return item, item != nil
		},
	}
}

func (gui *Gui) stashListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "stash",
			WindowName: "stash",
			Key:        STASH_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:             func() int { return len(gui.State.StashEntries) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Stash },
		OnFocus:                    gui.handleStashEntrySelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		GetDisplayStrings: func() [][]string {
			return presentation.GetStashEntryListDisplayStrings(gui.State.StashEntries, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedStashEntry()
			return item, item != nil
		},
	}
}

func (gui *Gui) commitFilesListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "commitFiles",
			WindowName: "commits",
			Key:        COMMIT_FILES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:             func() int { return gui.State.CommitFileManager.GetItemsLength() },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.CommitFiles },
		OnFocus:                    gui.handleCommitFileSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		GetDisplayStrings: func() [][]string {
			if gui.State.CommitFileManager.GetItemsLength() == 0 {
				return [][]string{{utils.ColoredString("(none)", color.FgRed)}}
			}

			lines := gui.State.CommitFileManager.Render(gui.State.Modes.Diffing.Ref, gui.State.Modes.PatchManager)
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

func (gui *Gui) submodulesListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "files",
			WindowName: "files",
			Key:        SUBMODULES_CONTEXT_KEY,
			Kind:       SIDE_CONTEXT,
		},
		GetItemsLength:             func() int { return len(gui.State.Submodules) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Submodules },
		OnFocus:                    gui.handleSubmoduleSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		GetDisplayStrings: func() [][]string {
			return presentation.GetSubmoduleListDisplayStrings(gui.State.Submodules)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedSubmodule()
			return item, item != nil
		},
	}
}

func (gui *Gui) suggestionsListContext() *ListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "suggestions",
			WindowName: "suggestions",
			Key:        SUGGESTIONS_CONTEXT_KEY,
			Kind:       PERSISTENT_POPUP,
		},
		GetItemsLength:             func() int { return len(gui.State.Suggestions) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Suggestions },
		OnFocus:                    func() error { return nil },
		Gui:                        gui,
		ResetMainViewOriginOnFocus: false,
		GetDisplayStrings: func() [][]string {
			return presentation.GetSuggestionListDisplayStrings(gui.State.Suggestions)
		},
	}
}

func (gui *Gui) getListContexts() []*ListContext {
	return []*ListContext{
		gui.State.Contexts.Menu,
		gui.State.Contexts.Files,
		gui.State.Contexts.Branches,
		gui.State.Contexts.Remotes,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Tags,
		gui.State.Contexts.BranchCommits,
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
			{ViewName: listContext.ViewName, Tag: "navigation", Contexts: []string{string(listContext.Key)}, Key: gui.getKey(keybindingConfig.Universal.PrevItemAlt), Modifier: gocui.ModNone, Handler: listContext.handlePrevLine},
			{ViewName: listContext.ViewName, Tag: "navigation", Contexts: []string{string(listContext.Key)}, Key: gui.getKey(keybindingConfig.Universal.PrevItem), Modifier: gocui.ModNone, Handler: listContext.handlePrevLine},
			{ViewName: listContext.ViewName, Tag: "navigation", Contexts: []string{string(listContext.Key)}, Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: listContext.handlePrevLine},
			{ViewName: listContext.ViewName, Tag: "navigation", Contexts: []string{string(listContext.Key)}, Key: gui.getKey(keybindingConfig.Universal.NextItemAlt), Modifier: gocui.ModNone, Handler: listContext.handleNextLine},
			{ViewName: listContext.ViewName, Tag: "navigation", Contexts: []string{string(listContext.Key)}, Key: gui.getKey(keybindingConfig.Universal.NextItem), Modifier: gocui.ModNone, Handler: listContext.handleNextLine},
			{ViewName: listContext.ViewName, Tag: "navigation", Contexts: []string{string(listContext.Key)}, Key: gui.getKey(keybindingConfig.Universal.PrevPage), Modifier: gocui.ModNone, Handler: listContext.handlePrevPage, Description: gui.Tr.LcPrevPage},
			{ViewName: listContext.ViewName, Tag: "navigation", Contexts: []string{string(listContext.Key)}, Key: gui.getKey(keybindingConfig.Universal.NextPage), Modifier: gocui.ModNone, Handler: listContext.handleNextPage, Description: gui.Tr.LcNextPage},
			{ViewName: listContext.ViewName, Tag: "navigation", Contexts: []string{string(listContext.Key)}, Key: gui.getKey(keybindingConfig.Universal.GotoTop), Modifier: gocui.ModNone, Handler: listContext.handleGotoTop, Description: gui.Tr.LcGotoTop},
			{ViewName: listContext.ViewName, Tag: "navigation", Contexts: []string{string(listContext.Key)}, Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: listContext.handleNextLine},
			{ViewName: listContext.ViewName, Contexts: []string{string(listContext.Key)}, Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: listContext.handleClick},
		}...)

		// the commits panel needs to lazyload things so it has a couple of its own handlers
		openSearchHandler := gui.handleOpenSearch
		gotoBottomHandler := listContext.handleGotoBottom
		if listContext.ViewName == "commits" {
			openSearchHandler = gui.handleOpenSearchForCommitsPanel
			gotoBottomHandler = gui.handleGotoBottomForCommitsPanel
		}

		bindings = append(bindings, []*Binding{
			{
				ViewName:    listContext.ViewName,
				Contexts:    []string{string(listContext.Key)},
				Key:         gui.getKey(keybindingConfig.Universal.StartSearch),
				Handler:     func() error { return openSearchHandler(listContext.ViewName) },
				Description: gui.Tr.LcStartSearch,
				Tag:         "navigation",
			},
			{
				ViewName:    listContext.ViewName,
				Contexts:    []string{string(listContext.Key)},
				Key:         gui.getKey(keybindingConfig.Universal.GotoBottom),
				Handler:     gotoBottomHandler,
				Description: gui.Tr.LcGotoBottom,
				Tag:         "navigation",
			},
		}...)
	}

	return bindings
}
