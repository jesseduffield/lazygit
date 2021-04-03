package gui

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

const SEARCH_PREFIX = "search: "
const INFO_SECTION_PADDING = " "

func (gui *Gui) informationStr() string {
	for _, mode := range gui.modeStatuses() {
		if mode.isActive() {
			return mode.description()
		}
	}

	if gui.g.Mouse {
		donate := color.New(color.FgMagenta, color.Underline).Sprint(gui.Tr.Donate)
		return donate + " " + gui.Config.GetVersion()
	} else {
		return gui.Config.GetVersion()
	}
}

// layout is called for every screen re-render e.g. when the screen is resized
func (gui *Gui) layout(g *gocui.Gui) error {
	g.Highlight = true
	width, height := g.Size()

	minimumHeight := 9
	minimumWidth := 10
	if height < minimumHeight || width < minimumWidth {
		v, err := g.SetView("limit", 0, 0, width-1, height-1, 0)
		if err != nil {
			if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
				return err
			}
			v.Title = gui.Tr.NotEnoughSpace
			v.Wrap = true
			_, _ = g.SetViewOnTop("limit")
		}
		return nil
	}

	informationStr := gui.informationStr()
	appStatus := gui.statusManager.getStatusString()

	viewDimensions := gui.getWindowDimensions(informationStr, appStatus)

	_, _ = g.SetViewOnBottom("limit")
	_ = g.DeleteView("limit")

	textColor := theme.GocuiDefaultTextColor

	// reading more lines into main view buffers upon resize
	prevMainView, err := gui.g.View("main")
	if err == nil {
		_, prevMainHeight := prevMainView.Size()
		newMainHeight := viewDimensions["main"].Y1 - viewDimensions["main"].Y0 - 1
		heightDiff := newMainHeight - prevMainHeight
		if heightDiff > 0 {
			if manager, ok := gui.viewBufferManagerMap["main"]; ok {
				manager.ReadLines(heightDiff)
			}
			if manager, ok := gui.viewBufferManagerMap["secondary"]; ok {
				manager.ReadLines(heightDiff)
			}
		}
	}

	setViewFromDimensions := func(viewName string, windowName string, frame bool) (*gocui.View, error) {
		dimensionsObj, ok := viewDimensions[windowName]

		if !ok {
			// view not specified in dimensions object: so create the view and hide it
			// making the view take up the whole space in the background in case it needs
			// to render content as soon as it appears, because lazyloaded content (via a pty task)
			// cares about the size of the view.
			view, err := g.SetView(viewName, 0, 0, width, height, 0)
			if err != nil {
				return view, err
			}
			return g.SetViewOnBottom(viewName)
		}

		frameOffset := 1
		if frame {
			frameOffset = 0
		}
		return g.SetView(
			viewName,
			dimensionsObj.X0-frameOffset,
			dimensionsObj.Y0-frameOffset,
			dimensionsObj.X1+frameOffset,
			dimensionsObj.Y1+frameOffset,
			0,
		)
	}

	v, err := setViewFromDimensions("main", "main", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		v.Title = gui.Tr.DiffTitle
		v.Wrap = true
		v.FgColor = textColor
		v.IgnoreCarriageReturns = true
	}

	secondaryView, err := setViewFromDimensions("secondary", "secondary", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		secondaryView.Title = gui.Tr.DiffTitle
		secondaryView.Wrap = true
		secondaryView.FgColor = textColor
		secondaryView.IgnoreCarriageReturns = true
	}

	hiddenViewOffset := 9999

	if v, err := setViewFromDimensions("status", "status", true); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		v.Title = gui.Tr.StatusTitle
		v.FgColor = textColor
	}

	filesView, err := setViewFromDimensions("files", "files", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		filesView.Highlight = true
		filesView.Title = gui.Tr.FilesTitle
		filesView.FgColor = textColor
		filesView.ContainsList = true
	}

	branchesView, err := setViewFromDimensions("branches", "branches", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		branchesView.Title = gui.Tr.BranchesTitle
		branchesView.FgColor = textColor
		branchesView.ContainsList = true
	}

	commitFilesView, err := setViewFromDimensions("commitFiles", gui.State.Contexts.CommitFiles.GetWindowName(), true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		commitFilesView.Title = gui.Tr.CommitFiles
		commitFilesView.FgColor = textColor
		commitFilesView.ContainsList = true
		_, _ = gui.g.SetViewOnBottom("commitFiles")
	}

	commitsView, err := setViewFromDimensions("commits", "commits", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		commitsView.Title = gui.Tr.CommitsTitle
		commitsView.FgColor = textColor
		commitsView.ContainsList = true
	}

	stashView, err := setViewFromDimensions("stash", "stash", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		stashView.Title = gui.Tr.StashTitle
		stashView.FgColor = textColor
		stashView.ContainsList = true
	}

	if gui.getCommitMessageView() == nil {
		// doesn't matter where this view starts because it will be hidden
		if commitMessageView, err := g.SetView("commitMessage", hiddenViewOffset, hiddenViewOffset, hiddenViewOffset+10, hiddenViewOffset+10, 0); err != nil {
			if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
				return err
			}
			_, _ = g.SetViewOnBottom("commitMessage")
			commitMessageView.Title = gui.Tr.CommitMessage
			commitMessageView.FgColor = textColor
			commitMessageView.Editable = true
			commitMessageView.Editor = gocui.EditorFunc(gui.commitMessageEditor)
		}
	}

	if check, _ := g.View("credentials"); check == nil {
		// doesn't matter where this view starts because it will be hidden
		if credentialsView, err := g.SetView("credentials", hiddenViewOffset, hiddenViewOffset, hiddenViewOffset+10, hiddenViewOffset+10, 0); err != nil {
			if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
				return err
			}
			_, _ = g.SetViewOnBottom("credentials")
			credentialsView.Title = gui.Tr.CredentialsUsername
			credentialsView.FgColor = textColor
			credentialsView.Editable = true
		}
	}

	if v, err := setViewFromDimensions("options", "options", false); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		v.Frame = false
		v.FgColor = theme.OptionsColor
	}

	// this view takes up one character. Its only purpose is to show the slash when searching
	if searchPrefixView, err := setViewFromDimensions("searchPrefix", "searchPrefix", false); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}

		searchPrefixView.BgColor = gocui.ColorDefault
		searchPrefixView.FgColor = gocui.ColorGreen
		searchPrefixView.Frame = false
		gui.setViewContent(searchPrefixView, SEARCH_PREFIX)
	}

	if searchView, err := setViewFromDimensions("search", "search", false); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}

		searchView.BgColor = gocui.ColorDefault
		searchView.FgColor = gocui.ColorGreen
		searchView.Frame = false
		searchView.Editable = true
	}

	if appStatusView, err := setViewFromDimensions("appStatus", "appStatus", false); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		appStatusView.BgColor = gocui.ColorDefault
		appStatusView.FgColor = gocui.ColorCyan
		appStatusView.Frame = false
		_, _ = g.SetViewOnBottom("appStatus")
	}

	informationView, err := setViewFromDimensions("information", "information", false)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		informationView.BgColor = gocui.ColorDefault
		informationView.FgColor = gocui.ColorGreen
		informationView.Frame = false
		gui.renderString("information", INFO_SECTION_PADDING+informationStr)
	}
	if gui.State.OldInformation != informationStr {
		gui.setViewContent(informationView, informationStr)
		gui.State.OldInformation = informationStr
	}

	if !gui.ViewsSetup {
		if err := gui.onInitialViewsCreation(); err != nil {
			return err
		}

		gui.ViewsSetup = true
	}

	if !gui.State.ViewsSetup {
		if err := gui.onInitialViewsCreationForRepo(); err != nil {
			return err
		}

		gui.State.ViewsSetup = true
	}

	for _, listContext := range gui.getListContexts() {
		view, err := gui.g.View(listContext.ViewName)
		if err != nil {
			continue
		}

		// ignore contexts whose view is owned by another context right now
		if view.Context != listContext.GetKey() {
			continue
		}

		// check if the selected line is now out of view and if so refocus it
		view.FocusPoint(0, listContext.GetPanelState().GetSelectedLineIdx())

		view.SelBgColor = theme.GocuiSelectedLineBgColor

		// I doubt this is expensive though it's admittedly redundant after the first render
		view.SetOnSelectItem(gui.onSelectItemWrapper(listContext.onSearchSelect))
	}

	gui.getMainView().SetOnSelectItem(gui.onSelectItemWrapper(gui.handlelineByLineNavigateTo))

	mainViewWidth, mainViewHeight := gui.getMainView().Size()
	if mainViewWidth != gui.State.PrevMainWidth || mainViewHeight != gui.State.PrevMainHeight {
		gui.State.PrevMainWidth = mainViewWidth
		gui.State.PrevMainHeight = mainViewHeight
		if err := gui.onResize(); err != nil {
			return err
		}
	}

	// here is a good place log some stuff
	// if you run `lazygit --logs`
	// this will let you see these branches as prettified json
	// gui.Log.Info(utils.AsJson(gui.State.Branches[0:4]))
	return gui.resizeCurrentPopupPanel()
}

func (gui *Gui) onInitialViewsCreationForRepo() error {
	gui.setInitialViewContexts()

	// hide any popup views. This only applies when we've just switched repos
	for _, viewName := range gui.popupViewNames() {
		_, _ = gui.g.SetViewOnBottom(viewName)
	}

	// the status panel is not actually a list context at the moment, so it is excluded
	// here. Arguably that's quite convenient because it means we're back to starting
	// in the files panel when landing in a new repo, but when returning from a submodule
	// we'll be back in the submodules context. This still seems awkward though, and it's
	// definitely going to break when (if) we make the status context a list context
	initialContext := gui.currentSideContext()
	if initialContext == nil {
		if gui.State.Modes.Filtering.Active() {
			initialContext = gui.State.Contexts.BranchCommits
		} else {
			initialContext = gui.State.Contexts.Files
		}
	}

	if err := gui.pushContext(initialContext); err != nil {
		return err
	}

	return gui.loadNewRepo()
}

func (gui *Gui) onInitialViewsCreation() error {
	gui.g.Mutexes.ViewsMutex.Lock()
	// add tabs to views
	for _, view := range gui.g.Views() {
		tabs := gui.viewTabNames(view.Name())
		if len(tabs) == 0 {
			continue
		}
		view.Tabs = tabs
	}
	gui.g.Mutexes.ViewsMutex.Unlock()

	if err := gui.keybindings(); err != nil {
		return err
	}

	if gui.showRecentRepos {
		if err := gui.handleCreateRecentReposMenu(); err != nil {
			return err
		}
		gui.showRecentRepos = false
	}

	gui.Updater.CheckForNewUpdate(gui.onBackgroundUpdateCheckFinish, false)

	gui.waitForIntro.Done()

	return nil
}
