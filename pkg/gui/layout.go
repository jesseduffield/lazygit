package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

const SEARCH_PREFIX = "search: "

// layout is called for every screen re-render e.g. when the screen is resized
func (gui *Gui) layout(g *gocui.Gui) error {
	if !gui.ViewsSetup {
		gui.printCommandLogHeader()

		if _, err := gui.g.SetCurrentView(gui.defaultSideContext().GetViewName()); err != nil {
			return err
		}
	}

	g.Highlight = true
	width, height := g.Size()

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
			view.Frame = frame
			view.Visible = true
		}

		return view, err
	}

	for _, arg := range gui.controlledViews() {
		_, err := setViewFromDimensions(arg.viewName, arg.windowName, arg.frame)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
	}

	minimumHeight := 9
	minimumWidth := 10
	gui.Views.Limit.Visible = height < minimumHeight || width < minimumWidth

	gui.Views.Tooltip.Visible = gui.Views.Menu.Visible && gui.Views.Tooltip.Buffer() != ""

	for _, context := range gui.TransientContexts() {
		view, err := gui.g.View(context.GetViewName())
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		view.Visible = gui.getViewNameForWindow(context.GetWindowName()) == context.GetViewName()
	}

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

		if !gui.isContextVisible(listContext) {
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
	// gui.c.Log.Info(utils.AsJson(gui.State.Model.Branches[0:4]))
	return gui.resizeCurrentPopupPanel()
}

func (gui *Gui) prepareView(viewName string) (*gocui.View, error) {
	// arbitrarily giving the view enough size so that we don't get an error, but
	// it's expected that the view will be given the correct size before being shown
	return gui.g.SetView(viewName, 0, 0, 10, 10, 0)
}

func (gui *Gui) onInitialViewsCreationForRepo() error {
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
	for _, view := range gui.orderedViews() {
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
