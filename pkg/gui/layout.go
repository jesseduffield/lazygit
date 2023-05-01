package gui

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
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

	// we assume that the view has already been created.
	setViewFromDimensions := func(viewName string, windowName string) (*gocui.View, error) {
		dimensionsObj, ok := viewDimensions[windowName]

		view, err := g.View(viewName)
		if err != nil {
			return nil, err
		}

		if !ok {
			// view not specified in dimensions object: so create the view and hide it
			// making the view take up the whole space in the background in case it needs
			// to render content as soon as it appears, because lazyloaded content (via a pty task)
			// cares about the size of the view.
			_, err := g.SetView(viewName, 0, 0, width, height, 0)
			view.Visible = false
			return view, err
		}

		frameOffset := 1
		if view.Frame {
			frameOffset = 0
		}
		_, err = g.SetView(
			viewName,
			dimensionsObj.X0-frameOffset,
			dimensionsObj.Y0-frameOffset,
			dimensionsObj.X1+frameOffset,
			dimensionsObj.Y1+frameOffset,
			0,
		)
		view.Visible = true

		return view, err
	}

	for _, context := range gui.State.Contexts.Flatten() {
		if !context.HasControlledBounds() {
			continue
		}

		_, err := setViewFromDimensions(context.GetViewName(), context.GetWindowName())
		if err != nil && !gocui.IsUnknownView(err) {
			return err
		}
	}

	minimumHeight := 9
	minimumWidth := 10
	gui.Views.Limit.Visible = height < minimumHeight || width < minimumWidth

	gui.Views.Tooltip.Visible = gui.Views.Menu.Visible && gui.Views.Tooltip.Buffer() != ""

	for _, context := range gui.TransientContexts() {
		view, err := gui.g.View(context.GetViewName())
		if err != nil && !gocui.IsUnknownView(err) {
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

		listContext.FocusLine()

		view.SelBgColor = theme.GocuiSelectedLineBgColor

		// I doubt this is expensive though it's admittedly redundant after the first render
		view.SetOnSelectItem(gui.onSelectItemWrapper(listContext.OnSearchSelect))
	}

	for _, context := range gui.getPatchExplorerContexts() {
		context := context
		context.GetView().SetOnSelectItem(gui.onSelectItemWrapper(
			func(selectedLineIdx int) error {
				context.GetMutex().Lock()
				defer context.GetMutex().Unlock()
				return context.NavigateTo(gui.c.IsCurrentContext(context), selectedLineIdx)
			}),
		)
	}

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
	if err := gui.c.ActivateContext(initialContext); err != nil {
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
		// if the view is in our mapping, we'll set the tabs and the tab index
		for _, values := range gui.viewTabMap() {
			index := slices.IndexFunc(values, func(tabContext context.TabView) bool {
				return tabContext.ViewName == view.Name()
			})

			if index != -1 {
				view.Tabs = slices.Map(values, func(tabContext context.TabView) string {
					return tabContext.Tab
				})
				view.TabIndex = index
			}
		}
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
