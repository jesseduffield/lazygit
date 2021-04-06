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
	var err error
	gui.Views.Limit, err = g.SetView("limit", 0, 0, width-1, height-1, 0)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Limit.Title = gui.Tr.NotEnoughSpace
		gui.Views.Limit.Wrap = true
	}
	gui.Views.Limit.Visible = height < minimumHeight || width < minimumWidth

	informationStr := gui.informationStr()
	appStatus := gui.statusManager.getStatusString()

	viewDimensions := gui.getWindowDimensions(informationStr, appStatus)

	// reading more lines into main view buffers upon resize
	prevMainView := gui.Views.Main
	if prevMainView != nil {
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
			if view != nil {
				view.Visible = false
			}
			return view, err
		}

		frameOffset := 1
		if frame {
			frameOffset = 0
		}
		view, err := g.SetView(
			viewName,
			dimensionsObj.X0-frameOffset,
			dimensionsObj.Y0-frameOffset,
			dimensionsObj.X1+frameOffset,
			dimensionsObj.Y1+frameOffset,
			0,
		)

		if view != nil {
			view.Visible = true
		}

		return view, err
	}

	gui.Views.Main, err = setViewFromDimensions("main", "main", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Main.Title = gui.Tr.DiffTitle
		gui.Views.Main.Wrap = true
		gui.Views.Main.FgColor = theme.GocuiDefaultTextColor
		gui.Views.Main.IgnoreCarriageReturns = true
	}

	gui.Views.Secondary, err = setViewFromDimensions("secondary", "secondary", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Secondary.Title = gui.Tr.DiffTitle
		gui.Views.Secondary.Wrap = true
		gui.Views.Secondary.FgColor = theme.GocuiDefaultTextColor
		gui.Views.Secondary.IgnoreCarriageReturns = true
	}

	if gui.Views.Status, err = setViewFromDimensions("status", "status", true); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Status.Title = gui.Tr.StatusTitle
		gui.Views.Status.FgColor = theme.GocuiDefaultTextColor
	}

	gui.Views.Files, err = setViewFromDimensions("files", "files", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Files.Highlight = true
		gui.Views.Files.Title = gui.Tr.FilesTitle
		gui.Views.Files.FgColor = theme.GocuiDefaultTextColor
		gui.Views.Files.ContainsList = true
	}

	gui.Views.Branches, err = setViewFromDimensions("branches", "branches", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Branches.Title = gui.Tr.BranchesTitle
		gui.Views.Branches.FgColor = theme.GocuiDefaultTextColor
		gui.Views.Branches.ContainsList = true
	}

	gui.Views.CommitFiles, err = setViewFromDimensions("commitFiles", gui.State.Contexts.CommitFiles.GetWindowName(), true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.CommitFiles.Title = gui.Tr.CommitFiles
		gui.Views.CommitFiles.FgColor = theme.GocuiDefaultTextColor
		gui.Views.CommitFiles.ContainsList = true
	}
	// if the commit files view is the view to be displayed for its window, we'll display it
	gui.Views.CommitFiles.Visible = gui.getViewNameForWindow(gui.State.Contexts.CommitFiles.GetWindowName()) == "commitFiles"

	gui.Views.Commits, err = setViewFromDimensions("commits", "commits", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Commits.Title = gui.Tr.CommitsTitle
		gui.Views.Commits.FgColor = theme.GocuiDefaultTextColor
		gui.Views.Commits.ContainsList = true
	}

	gui.Views.Stash, err = setViewFromDimensions("stash", "stash", true)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Stash.Title = gui.Tr.StashTitle
		gui.Views.Stash.FgColor = theme.GocuiDefaultTextColor
		gui.Views.Stash.ContainsList = true
	}

	if gui.Views.Options, err = setViewFromDimensions("options", "options", false); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Options.Frame = false
		gui.Views.Options.FgColor = theme.OptionsColor
	}

	// this view takes up one character. Its only purpose is to show the slash when searching
	if gui.Views.SearchPrefix, err = setViewFromDimensions("searchPrefix", "searchPrefix", false); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}

		gui.Views.SearchPrefix.BgColor = gocui.ColorDefault
		gui.Views.SearchPrefix.FgColor = gocui.ColorGreen
		gui.Views.SearchPrefix.Frame = false
		gui.setViewContent(gui.Views.SearchPrefix, SEARCH_PREFIX)
	}

	if gui.Views.Search, err = setViewFromDimensions("search", "search", false); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}

		gui.Views.Search.BgColor = gocui.ColorDefault
		gui.Views.Search.FgColor = gocui.ColorGreen
		gui.Views.Search.Frame = false
		gui.Views.Search.Editable = true
	}

	if gui.Views.AppStatus, err = setViewFromDimensions("appStatus", "appStatus", false); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.AppStatus.BgColor = gocui.ColorDefault
		gui.Views.AppStatus.FgColor = gocui.ColorCyan
		gui.Views.AppStatus.Frame = false
		gui.Views.AppStatus.Visible = false
	}

	gui.Views.Information, err = setViewFromDimensions("information", "information", false)
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Information.BgColor = gocui.ColorDefault
		gui.Views.Information.FgColor = gocui.ColorGreen
		gui.Views.Information.Frame = false
		gui.renderString(gui.Views.Information, INFO_SECTION_PADDING+informationStr)
	}
	if gui.State.OldInformation != informationStr {
		gui.setViewContent(gui.Views.Information, informationStr)
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
		if ContextKey(view.Context) != listContext.GetKey() {
			continue
		}

		// check if the selected line is now out of view and if so refocus it
		view.FocusPoint(0, listContext.GetPanelState().GetSelectedLineIdx())

		view.SelBgColor = theme.GocuiSelectedLineBgColor

		// I doubt this is expensive though it's admittedly redundant after the first render
		view.SetOnSelectItem(gui.onSelectItemWrapper(listContext.onSearchSelect))
	}

	gui.Views.Main.SetOnSelectItem(gui.onSelectItemWrapper(gui.handlelineByLineNavigateTo))

	mainViewWidth, mainViewHeight := gui.Views.Main.Size()
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

func (gui *Gui) setHiddenView(viewName string) (*gocui.View, error) {
	// arbitrarily giving the view enough size so that we don't get an error, but
	// it's expected that the view will be given the correct size before being shown
	return gui.g.SetView(viewName, 0, 0, 10, 10, 0)
}

func (gui *Gui) onInitialViewsCreationForRepo() error {
	gui.setInitialViewContexts()

	// hide any popup views. This only applies when we've just switched repos
	for _, viewName := range gui.popupViewNames() {
		view, err := gui.g.View(viewName)
		if err == nil {
			view.Visible = false
		}
	}

	initialContext := gui.currentSideContext()
	if err := gui.pushContext(initialContext); err != nil {
		return err
	}

	return gui.loadNewRepo()
}

func (gui *Gui) onInitialViewsCreation() error {
	// creating some views which are hidden at the start but we need to exist so that we can set an initial ordering
	if err := gui.createHiddenViews(); err != nil {
		return err
	}

	// now we order the views (in order of bottom first)
	layerOneViews := []*gocui.View{
		// first layer. Ordering within this layer does not matter because there are
		// no overlapping views
		gui.Views.Status,
		gui.Views.Files,
		gui.Views.Branches,
		gui.Views.Commits,
		gui.Views.Stash,
		gui.Views.CommitFiles,
		gui.Views.Main,
		gui.Views.Secondary,

		// bottom line
		gui.Views.Options,
		gui.Views.AppStatus,
		gui.Views.Information,
		gui.Views.Search,
		gui.Views.SearchPrefix,

		// popups. Ordering within this layer does not matter because there should
		// only be one popup shown at a time
		gui.Views.CommitMessage,
		gui.Views.Credentials,
		gui.Views.Menu,
		gui.Views.Suggestions,
		gui.Views.Confirmation,

		// this guy will cover everything else when it appears
		gui.Views.Limit,
	}

	for _, view := range layerOneViews {
		if _, err := gui.g.SetViewOnTop(view.Name()); err != nil {
			return err
		}
	}

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

func (gui *Gui) createHiddenViews() error {
	// doesn't matter where this view starts because it will be hidden
	var err error
	if gui.Views.CommitMessage, err = gui.setHiddenView("commitMessage"); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.CommitMessage.Visible = false
		gui.Views.CommitMessage.Title = gui.Tr.CommitMessage
		gui.Views.CommitMessage.FgColor = theme.GocuiDefaultTextColor
		gui.Views.CommitMessage.Editable = true
		gui.Views.CommitMessage.Editor = gocui.EditorFunc(gui.commitMessageEditor)
	}

	// doesn't matter where this view starts because it will be hidden
	if gui.Views.Credentials, err = gui.setHiddenView("credentials"); err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Credentials.Visible = false
		gui.Views.Credentials.Title = gui.Tr.CredentialsUsername
		gui.Views.Credentials.FgColor = theme.GocuiDefaultTextColor
		gui.Views.Credentials.Editable = true
	}

	// not worrying about setting attributes because that will be done when the view is actually shown
	gui.Views.Confirmation, err = gui.setHiddenView("confirmation")
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Confirmation.Visible = false
	}

	// not worrying about setting attributes because that will be done when the view is actually shown
	gui.Views.Suggestions, err = gui.setHiddenView("suggestions")
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Suggestions.Visible = false
	}

	// not worrying about setting attributes because that will be done when the view is actually shown
	gui.Views.Menu, err = gui.setHiddenView("menu")
	if err != nil {
		if err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		gui.Views.Menu.Visible = false
	}

	return nil
}
