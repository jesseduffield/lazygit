package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

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

	appStatus := gui.helpers.AppStatus.GetStatusString()

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

	contextsToRerender := []types.Context{}

	// we assume that the view has already been created.
	setViewFromDimensions := func(context types.Context) (*gocui.View, error) {
		viewName := context.GetViewName()
		windowName := context.GetWindowName()

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

		if context.NeedsRerenderOnWidthChange() {
			// view.Width() returns the width -1 for some reason
			oldWidth := view.Width() + 1
			newWidth := dimensionsObj.X1 - dimensionsObj.X0 + 2*frameOffset
			if oldWidth != newWidth {
				contextsToRerender = append(contextsToRerender, context)
			}
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

		_, err := setViewFromDimensions(context)
		if err != nil && !gocui.IsUnknownView(err) {
			return err
		}
	}

	minimumHeight := 9
	minimumWidth := 10
	gui.Views.Limit.Visible = height < minimumHeight || width < minimumWidth

	gui.Views.Tooltip.Visible = gui.Views.Menu.Visible && gui.Views.Tooltip.Buffer() != ""

	for _, context := range gui.transientContexts() {
		view, err := gui.g.View(context.GetViewName())
		if err != nil && !gocui.IsUnknownView(err) {
			return err
		}
		view.Visible = gui.helpers.Window.GetViewNameForWindow(context.GetWindowName()) == context.GetViewName()
	}

	if gui.PrevLayout.Information != informationStr {
		gui.c.SetViewContent(gui.Views.Information, informationStr)
		gui.PrevLayout.Information = informationStr
	}

	if !gui.ViewsSetup {
		if err := gui.onInitialViewsCreation(); err != nil {
			return err
		}

		gui.handleTestMode()

		gui.ViewsSetup = true
	}

	if !gui.State.ViewsSetup {
		if err := gui.onInitialViewsCreationForRepo(); err != nil {
			return err
		}

		gui.State.ViewsSetup = true
	}

	mainViewWidth, mainViewHeight := gui.Views.Main.Size()
	if mainViewWidth != gui.PrevLayout.MainWidth || mainViewHeight != gui.PrevLayout.MainHeight {
		gui.PrevLayout.MainWidth = mainViewWidth
		gui.PrevLayout.MainHeight = mainViewHeight
		if err := gui.onResize(); err != nil {
			return err
		}
	}

	for _, context := range contextsToRerender {
		if err := context.HandleRender(); err != nil {
			return err
		}
	}

	// here is a good place log some stuff
	// if you run `lazygit --logs`
	// this will let you see these branches as prettified json
	// gui.c.Log.Info(utils.AsJson(gui.State.Model.Branches[0:4]))
	if err := gui.helpers.Confirmation.ResizeCurrentPopupPanel(); err != nil {
		return err
	}

	gui.renderContextOptionsMap()

outer:
	for {
		select {
		case f := <-gui.afterLayoutFuncs:
			if err := f(); err != nil {
				return err
			}
		default:
			break outer
		}
	}

	return nil
}

func (gui *Gui) prepareView(viewName string) (*gocui.View, error) {
	// arbitrarily giving the view enough size so that we don't get an error, but
	// it's expected that the view will be given the correct size before being shown
	return gui.g.SetView(viewName, 0, 0, 10, 10, 0)
}

func (gui *Gui) onInitialViewsCreationForRepo() error {
	if err := gui.onRepoViewReset(); err != nil {
		return err
	}

	// hide any popup views. This only applies when we've just switched repos
	for _, viewName := range gui.popupViewNames() {
		view, err := gui.g.View(viewName)
		if err == nil {
			view.Visible = false
		}
	}

	initialContext := gui.c.CurrentContext()
	if err := gui.c.ActivateContext(initialContext); err != nil {
		return err
	}

	return gui.loadNewRepo()
}

func (gui *Gui) popupViewNames() []string {
	popups := lo.Filter(gui.State.Contexts.Flatten(), func(c types.Context, _ int) bool {
		return c.GetKind() == types.PERSISTENT_POPUP || c.GetKind() == types.TEMPORARY_POPUP
	})

	return lo.Map(popups, func(c types.Context, _ int) string {
		return c.GetViewName()
	})
}

func (gui *Gui) onRepoViewReset() error {
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
				view.Tabs = lo.Map(values, func(tabContext context.TabView, _ int) string {
					return tabContext.Tab
				})
				view.TabIndex = index
			}
		}
	}
	gui.g.Mutexes.ViewsMutex.Unlock()

	return nil
}

func (gui *Gui) onInitialViewsCreation() error {
	if !gui.c.UserConfig.DisableStartupPopups {
		storedPopupVersion := gui.c.GetAppState().StartupPopupVersion
		if storedPopupVersion < StartupPopupVersion {
			gui.showIntroPopupMessage()
		} else {
			gui.showBreakingChangesMessage()
		}
	}

	gui.c.GetAppState().LastVersion = gui.Config.GetVersion()
	gui.c.SaveAppStateAndLogError()

	if gui.showRecentRepos {
		if err := gui.helpers.Repos.CreateRecentReposMenu(); err != nil {
			return err
		}
		gui.showRecentRepos = false
	}

	gui.helpers.Update.CheckForUpdateInBackground()

	gui.waitForIntro.Done()

	return nil
}

// getFocusLayout returns a manager function for when view gain and lose focus
func (gui *Gui) getFocusLayout() func(g *gocui.Gui) error {
	var previousView *gocui.View
	return func(g *gocui.Gui) error {
		newView := gui.g.CurrentView()
		// for now we don't consider losing focus to a popup panel as actually losing focus
		if newView != previousView && !gui.helpers.Confirmation.IsPopupPanel(newView.Name()) {
			if err := gui.onViewFocusLost(previousView); err != nil {
				return err
			}

			previousView = newView
		}
		return nil
	}
}

func (gui *Gui) onViewFocusLost(oldView *gocui.View) error {
	if oldView == nil {
		return nil
	}

	oldView.Highlight = false

	_ = oldView.SetOriginX(0)

	return nil
}

func (gui *Gui) transientContexts() []types.Context {
	return lo.Filter(gui.State.Contexts.Flatten(), func(context types.Context, _ int) bool {
		return context.IsTransient()
	})
}
