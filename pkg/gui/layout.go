package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

const SEARCH_PREFIX = "search: "

func (gui *Gui) createAllViews() error {
	viewNameMappings := []struct {
		viewPtr **gocui.View
		name    string
	}{
		{viewPtr: &gui.Views.Status, name: "status"},
		{viewPtr: &gui.Views.Files, name: "files"},
		{viewPtr: &gui.Views.Branches, name: "branches"},
		{viewPtr: &gui.Views.Commits, name: "commits"},
		{viewPtr: &gui.Views.Stash, name: "stash"},
		{viewPtr: &gui.Views.CommitFiles, name: "commitFiles"},
		{viewPtr: &gui.Views.Main, name: "main"},
		{viewPtr: &gui.Views.Secondary, name: "secondary"},
		{viewPtr: &gui.Views.Options, name: "options"},
		{viewPtr: &gui.Views.AppStatus, name: "appStatus"},
		{viewPtr: &gui.Views.Information, name: "information"},
		{viewPtr: &gui.Views.Search, name: "search"},
		{viewPtr: &gui.Views.SearchPrefix, name: "searchPrefix"},
		{viewPtr: &gui.Views.CommitMessage, name: "commitMessage"},
		{viewPtr: &gui.Views.Credentials, name: "credentials"},
		{viewPtr: &gui.Views.Menu, name: "menu"},
		{viewPtr: &gui.Views.Suggestions, name: "suggestions"},
		{viewPtr: &gui.Views.Confirmation, name: "confirmation"},
		{viewPtr: &gui.Views.Limit, name: "limit"},
		{viewPtr: &gui.Views.Extras, name: "extras"},
	}

	var err error
	for _, mapping := range viewNameMappings {
		*mapping.viewPtr, err = gui.prepareView(mapping.name)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
	}

	gui.Views.Options.Frame = false
	gui.Views.Options.FgColor = theme.OptionsColor

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

	gui.Views.Branches.Title = gui.c.Tr.BranchesTitle
	gui.Views.Branches.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Files.Highlight = true
	gui.Views.Files.Title = gui.c.Tr.FilesTitle
	gui.Views.Files.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Secondary.Title = gui.c.Tr.DiffTitle
	gui.Views.Secondary.Wrap = true
	gui.Views.Secondary.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Secondary.IgnoreCarriageReturns = true

	gui.Views.Main.Title = gui.c.Tr.DiffTitle
	gui.Views.Main.Wrap = true
	gui.Views.Main.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Main.IgnoreCarriageReturns = true

	gui.Views.Limit.Title = gui.c.Tr.NotEnoughSpace
	gui.Views.Limit.Wrap = true

	gui.Views.Status.Title = gui.c.Tr.StatusTitle
	gui.Views.Status.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Search.BgColor = gocui.ColorDefault
	gui.Views.Search.FgColor = gocui.ColorGreen
	gui.Views.Search.Frame = false
	gui.Views.Search.Editable = true

	gui.Views.AppStatus.BgColor = gocui.ColorDefault
	gui.Views.AppStatus.FgColor = gocui.ColorCyan
	gui.Views.AppStatus.Frame = false
	gui.Views.AppStatus.Visible = false

	gui.Views.CommitMessage.Visible = false
	gui.Views.CommitMessage.Title = gui.c.Tr.CommitMessage
	gui.Views.CommitMessage.FgColor = theme.GocuiDefaultTextColor
	gui.Views.CommitMessage.Editable = true
	gui.Views.CommitMessage.Editor = gocui.EditorFunc(gui.commitMessageEditor)

	gui.Views.Confirmation.Visible = false

	gui.Views.Credentials.Visible = false
	gui.Views.Credentials.Title = gui.c.Tr.CredentialsUsername
	gui.Views.Credentials.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Credentials.Editable = true

	gui.Views.Suggestions.Visible = false

	gui.Views.Menu.Visible = false

	gui.Views.Information.BgColor = gocui.ColorDefault
	gui.Views.Information.FgColor = gocui.ColorGreen
	gui.Views.Information.Frame = false

	gui.Views.Extras.Title = gui.c.Tr.CommandLog
	gui.Views.Extras.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Extras.Autoscroll = true
	gui.Views.Extras.Wrap = true

	gui.printCommandLogHeader()

	if _, err := gui.g.SetCurrentView(gui.defaultSideContext().GetViewName()); err != nil {
		return err
	}

	return nil
}

// layout is called for every screen re-render e.g. when the screen is resized
func (gui *Gui) layout(g *gocui.Gui) error {
	if !gui.ViewsSetup {
		if err := gui.createAllViews(); err != nil {
			return err
		}
	}

	g.Highlight = true
	width, height := g.Size()

	minimumHeight := 9
	minimumWidth := 10
	var err error
	_, err = g.SetView("limit", 0, 0, width-1, height-1, 0)
	if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
		return err
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

	args := []struct {
		viewName   string
		windowName string
		frame      bool
	}{
		{viewName: "main", windowName: "main", frame: true},
		{viewName: "secondary", windowName: "secondary", frame: true},
		{viewName: "status", windowName: "status", frame: true},
		{viewName: "files", windowName: "files", frame: true},
		{viewName: "branches", windowName: "branches", frame: true},
		{viewName: "commitFiles", windowName: gui.State.Contexts.CommitFiles.GetWindowName(), frame: true},
		{viewName: "commits", windowName: "commits", frame: true},
		{viewName: "stash", windowName: "stash", frame: true},
		{viewName: "options", windowName: "options", frame: false},
		{viewName: "searchPrefix", windowName: "searchPrefix", frame: false},
		{viewName: "search", windowName: "search", frame: false},
		{viewName: "appStatus", windowName: "appStatus", frame: false},
		{viewName: "information", windowName: "information", frame: false},
		{viewName: "extras", windowName: "extras", frame: true},
	}

	for _, arg := range args {
		_, err = setViewFromDimensions(arg.viewName, arg.windowName, arg.frame)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
	}

	// if the commit files view is the view to be displayed for its window, we'll display it
	gui.Views.CommitFiles.Visible = gui.getViewNameForWindow(gui.State.Contexts.CommitFiles.GetWindowName()) == "commitFiles"

	if gui.PrevLayout.Information != informationStr {
		gui.setViewContent(gui.Views.Information, informationStr)
		gui.PrevLayout.Information = informationStr
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
		view, err := gui.g.View(listContext.GetViewName())
		if err != nil {
			continue
		}

		// ignore contexts whose view is owned by another context right now
		if types.ContextKey(view.Context) != listContext.GetKey() {
			continue
		}

		listContext.FocusLine()

		view.SelBgColor = theme.GocuiSelectedLineBgColor

		// I doubt this is expensive though it's admittedly redundant after the first render
		view.SetOnSelectItem(gui.onSelectItemWrapper(listContext.OnSearchSelect))
	}

	gui.Views.Main.SetOnSelectItem(gui.onSelectItemWrapper(gui.handlelineByLineNavigateTo))

	mainViewWidth, mainViewHeight := gui.Views.Main.Size()
	if mainViewWidth != gui.PrevLayout.MainWidth || mainViewHeight != gui.PrevLayout.MainHeight {
		gui.PrevLayout.MainWidth = mainViewWidth
		gui.PrevLayout.MainHeight = mainViewHeight
		if err := gui.onResize(); err != nil {
			return err
		}
	}

	// here is a good place log some stuff
	// if you run `lazygit --logs`
	// this will let you see these branches as prettified json
	// gui.c.Log.Info(utils.AsJson(gui.State.Branches[0:4]))
	return gui.resizeCurrentPopupPanel()
}

func (gui *Gui) prepareView(viewName string) (*gocui.View, error) {
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
	if err := gui.c.PushContext(initialContext); err != nil {
		return err
	}

	return gui.loadNewRepo()
}

func (gui *Gui) onInitialViewsCreation() error {
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
		gui.Views.Extras,

		// bottom line
		gui.Views.Options,
		gui.Views.AppStatus,
		gui.Views.Information,
		gui.Views.Search,
		gui.Views.SearchPrefix, // this view takes up one character. Its only purpose is to show the slash when searching

		// popups. Ordering within this layer does not matter because there should
		// only be one popup shown at a time
		gui.Views.CommitMessage,
		gui.Views.Menu,
		gui.Views.Suggestions,
		gui.Views.Confirmation,
		gui.Views.Credentials,

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

	if !gui.c.UserConfig.DisableStartupPopups {
		popupTasks := []func(chan struct{}) error{}
		storedPopupVersion := gui.c.GetAppState().StartupPopupVersion
		if storedPopupVersion < StartupPopupVersion {
			popupTasks = append(popupTasks, gui.showIntroPopupMessage)
		}
		gui.showInitialPopups(popupTasks)
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
