package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// getFocusLayout returns a manager function for when view gain and lose focus
func (gui *Gui) getFocusLayout() func(g *gocui.Gui) error {
	var previousView *gocui.View
	return func(g *gocui.Gui) error {
		newView := gui.g.CurrentView()
		if err := gui.onFocusChange(); err != nil {
			return err
		}
		// for now we don't consider losing focus to a popup panel as actually losing focus
		if newView != previousView && !gui.isPopupPanel(newView.Name()) {
			if err := gui.onFocusLost(previousView, newView); err != nil {
				return err
			}
			if err := gui.onFocus(newView); err != nil {
				return err
			}
			previousView = newView
		}
		return nil
	}
}

func (gui *Gui) onFocusChange() error {
	currentView := gui.g.CurrentView()
	for _, view := range gui.g.Views() {
		view.Highlight = view == currentView
	}
	return nil
}

func (gui *Gui) onFocusLost(v *gocui.View, newView *gocui.View) error {
	if v == nil {
		return nil
	}
	if v.IsSearching() && newView.Name() != "search" {
		if err := gui.onSearchEscape(); err != nil {
			return err
		}
	}
	switch v.Name() {
	case "main":
		// if we have lost focus to a first-class panel, we need to do some cleanup
		gui.changeMainViewsContext("normal")
	case "commitFiles":
		if gui.State.MainContext != "patch-building" {
			if _, err := gui.g.SetViewOnBottom(v.Name()); err != nil {
				return err
			}
		}
	}
	gui.Log.Info(v.Name() + " focus lost")
	return nil
}

func (gui *Gui) onFocus(v *gocui.View) error {
	if v == nil {
		return nil
	}
	gui.Log.Info(v.Name() + " focus gained")
	return nil
}

func (gui *Gui) getViewHeights() map[string]int {
	currView := gui.g.CurrentView()
	currentCyclebleView := gui.State.PreviousView
	if currView != nil {
		viewName := currView.Name()
		usePreviousView := true
		for _, view := range cyclableViews {
			if view == viewName {
				currentCyclebleView = viewName
				usePreviousView = false
				break
			}
		}
		if usePreviousView {
			currentCyclebleView = gui.State.PreviousView
		}
	}

	// unfortunate result of the fact that these are separate views, have to map explicitly
	if currentCyclebleView == "commitFiles" {
		currentCyclebleView = "commits"
	}

	_, height := gui.g.Size()

	if gui.State.ScreenMode == SCREEN_FULL || gui.State.ScreenMode == SCREEN_HALF {
		vHeights := map[string]int{
			"status":   0,
			"files":    0,
			"branches": 0,
			"commits":  0,
			"stash":    0,
			"options":  0,
		}
		vHeights[currentCyclebleView] = height - 1
		return vHeights
	}

	usableSpace := height - 7
	extraSpace := usableSpace - (usableSpace/3)*3

	if height >= 28 {
		return map[string]int{
			"status":   3,
			"files":    (usableSpace / 3) + extraSpace,
			"branches": usableSpace / 3,
			"commits":  usableSpace / 3,
			"stash":    3,
			"options":  1,
		}
	}

	defaultHeight := 3
	if height < 21 {
		defaultHeight = 1
	}
	vHeights := map[string]int{
		"status":   defaultHeight,
		"files":    defaultHeight,
		"branches": defaultHeight,
		"commits":  defaultHeight,
		"stash":    defaultHeight,
		"options":  defaultHeight,
	}
	vHeights[currentCyclebleView] = height - defaultHeight*4 - 1

	return vHeights
}

// layout is called for every screen re-render e.g. when the screen is resized
func (gui *Gui) layout(g *gocui.Gui) error {
	g.Highlight = true
	width, height := g.Size()

	information := gui.Config.GetVersion()
	if gui.g.Mouse {
		donate := color.New(color.FgMagenta, color.Underline).Sprint(gui.Tr.SLocalize("Donate"))
		information = donate + " " + information
	}
	if gui.inDiffMode() {
		information = utils.ColoredString(fmt.Sprintf("%s %s %s", gui.Tr.SLocalize("showingGitDiff"), "git diff "+gui.diffStr(), utils.ColoredString(gui.Tr.SLocalize("(reset)"), color.Underline)), color.FgMagenta)
	} else if gui.inFilterMode() {
		information = utils.ColoredString(fmt.Sprintf("%s '%s' %s", gui.Tr.SLocalize("filteringBy"), gui.State.FilterPath, utils.ColoredString(gui.Tr.SLocalize("(reset)"), color.Underline)), color.FgRed, color.Bold)
	} else if len(gui.State.CherryPickedCommits) > 0 {
		information = utils.ColoredString(fmt.Sprintf("%d commits copied", len(gui.State.CherryPickedCommits)), color.FgCyan)
	}

	minimumHeight := 9
	minimumWidth := 10
	if height < minimumHeight || width < minimumWidth {
		v, err := g.SetView("limit", 0, 0, width-1, height-1, 0)
		if err != nil {
			if err.Error() != "unknown view" {
				return err
			}
			v.Title = gui.Tr.SLocalize("NotEnoughSpace")
			v.Wrap = true
			_, _ = g.SetViewOnTop("limit")
		}
		return nil
	}

	vHeights := gui.getViewHeights()

	optionsVersionBoundary := width - max(len(utils.Decolorise(information)), 1)

	appStatus := gui.statusManager.getStatusString()
	appStatusOptionsBoundary := 0
	if appStatus != "" {
		appStatusOptionsBoundary = len(appStatus) + 2
	}

	_, _ = g.SetViewOnBottom("limit")
	_ = g.DeleteView("limit")

	sidePanelWidthRatio := gui.Config.GetUserConfig().GetFloat64("gui.sidePanelWidth")

	textColor := theme.GocuiDefaultTextColor
	var leftSideWidth int
	switch gui.State.ScreenMode {
	case SCREEN_NORMAL:
		leftSideWidth = int(float64(width) * sidePanelWidthRatio)
	case SCREEN_HALF:
		leftSideWidth = width/2 - 2
	case SCREEN_FULL:
		currentView := gui.g.CurrentView()
		if currentView != nil && currentView.Name() == "main" {
			leftSideWidth = 0
		} else {
			leftSideWidth = width - 1
		}
	}

	mainPanelLeft := leftSideWidth + 1
	mainPanelRight := width - 1
	secondaryPanelLeft := width - 1
	secondaryPanelTop := 0
	mainPanelBottom := height - 2
	if gui.State.SplitMainPanel {
		if gui.State.ScreenMode == SCREEN_FULL {
			mainPanelLeft = 0
			panelSplitX := width/2 - 4
			mainPanelRight = panelSplitX
			secondaryPanelLeft = panelSplitX + 1
		} else if width < 220 {
			mainPanelBottom = height/2 - 1
			secondaryPanelTop = mainPanelBottom + 1
			secondaryPanelLeft = leftSideWidth + 1
		} else {
			units := 5
			leftSideWidth = width / units
			mainPanelLeft = leftSideWidth + 1
			panelSplitX := (1 + ((units - 1) / 2)) * width / units
			mainPanelRight = panelSplitX
			secondaryPanelLeft = panelSplitX + 1
		}
	}

	main := "main"
	secondary := "secondary"
	swappingMainPanels := gui.State.Panels.LineByLine != nil && gui.State.Panels.LineByLine.SecondaryFocused
	if swappingMainPanels {
		main = "secondary"
		secondary = "main"
	}

	// reading more lines into main view buffers upon resize
	prevMainView, err := gui.g.View("main")
	if err == nil {
		_, prevMainHeight := prevMainView.Size()
		heightDiff := mainPanelBottom - prevMainHeight - 1
		if heightDiff > 0 {
			if manager, ok := gui.viewBufferManagerMap["main"]; ok {
				manager.ReadLines(heightDiff)
			}
			if manager, ok := gui.viewBufferManagerMap["secondary"]; ok {
				manager.ReadLines(heightDiff)
			}
		}
	}

	v, err := g.SetView(main, mainPanelLeft, 0, mainPanelRight, mainPanelBottom, gocui.LEFT)
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.Title = gui.Tr.SLocalize("DiffTitle")
		v.Wrap = true
		v.FgColor = textColor
		v.IgnoreCarriageReturns = true
	}

	hiddenViewOffset := 9999

	hiddenSecondaryPanelOffset := 0
	if !gui.State.SplitMainPanel {
		hiddenSecondaryPanelOffset = hiddenViewOffset
	}
	secondaryView, err := g.SetView(secondary, secondaryPanelLeft+hiddenSecondaryPanelOffset, hiddenSecondaryPanelOffset+secondaryPanelTop, width-1+hiddenSecondaryPanelOffset, height-2+hiddenSecondaryPanelOffset, gocui.LEFT)
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		secondaryView.Title = gui.Tr.SLocalize("DiffTitle")
		secondaryView.Wrap = true
		secondaryView.FgColor = gocui.ColorWhite
		secondaryView.IgnoreCarriageReturns = true
	}

	if v, err := g.SetView("status", 0, 0, leftSideWidth, vHeights["status"]-1, gocui.BOTTOM|gocui.RIGHT); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.Title = gui.Tr.SLocalize("StatusTitle")
		v.FgColor = textColor
	}

	filesView, err := g.SetViewBeneath("files", "status", vHeights["files"])
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		filesView.Highlight = true
		filesView.Title = gui.Tr.SLocalize("FilesTitle")
		filesView.SetOnSelectItem(gui.onSelectItemWrapper(gui.onFilesPanelSearchSelect))
		filesView.ContainsList = true
	}

	branchesView, err := g.SetViewBeneath("branches", "files", vHeights["branches"])
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		branchesView.Title = gui.Tr.SLocalize("BranchesTitle")
		branchesView.Tabs = []string{"Local Branches", "Remotes", "Tags"}
		branchesView.FgColor = textColor
		branchesView.SetOnSelectItem(gui.onSelectItemWrapper(gui.onBranchesPanelSearchSelect))
		branchesView.ContainsList = true
	}

	commitFilesView, err := g.SetViewBeneath("commitFiles", "branches", vHeights["commits"])
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		commitFilesView.Title = gui.Tr.SLocalize("CommitFiles")
		commitFilesView.FgColor = textColor
		commitFilesView.SetOnSelectItem(gui.onSelectItemWrapper(gui.onCommitFilesPanelSearchSelect))
		commitFilesView.ContainsList = true
	}

	commitsView, err := g.SetViewBeneath("commits", "branches", vHeights["commits"])
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		commitsView.Title = gui.Tr.SLocalize("CommitsTitle")
		commitsView.Tabs = []string{"Commits", "Reflog"}
		commitsView.FgColor = textColor
		commitsView.SetOnSelectItem(gui.onSelectItemWrapper(gui.onCommitsPanelSearchSelect))
		commitsView.ContainsList = true
	}

	stashView, err := g.SetViewBeneath("stash", "commits", vHeights["stash"])
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		stashView.Title = gui.Tr.SLocalize("StashTitle")
		stashView.FgColor = textColor
		stashView.SetOnSelectItem(gui.onSelectItemWrapper(gui.onStashPanelSearchSelect))
		stashView.ContainsList = true
	}

	if v, err := g.SetView("options", appStatusOptionsBoundary-1, height-2, optionsVersionBoundary-1, height, 0); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.Frame = false
		v.FgColor = theme.OptionsColor
	}

	if gui.getCommitMessageView() == nil {
		// doesn't matter where this view starts because it will be hidden
		if commitMessageView, err := g.SetView("commitMessage", hiddenViewOffset, hiddenViewOffset, hiddenViewOffset+10, hiddenViewOffset+10, 0); err != nil {
			if err.Error() != "unknown view" {
				return err
			}
			_, _ = g.SetViewOnBottom("commitMessage")
			commitMessageView.Title = gui.Tr.SLocalize("CommitMessage")
			commitMessageView.FgColor = textColor
			commitMessageView.Editable = true
			commitMessageView.Editor = gocui.EditorFunc(gui.commitMessageEditor)
		}
	}

	if check, _ := g.View("credentials"); check == nil {
		// doesn't matter where this view starts because it will be hidden
		if credentialsView, err := g.SetView("credentials", hiddenViewOffset, hiddenViewOffset, hiddenViewOffset+10, hiddenViewOffset+10, 0); err != nil {
			if err.Error() != "unknown view" {
				return err
			}
			_, err := g.SetViewOnBottom("credentials")
			if err != nil {
				return err
			}
			credentialsView.Title = gui.Tr.SLocalize("CredentialsUsername")
			credentialsView.FgColor = textColor
			credentialsView.Editable = true
		}
	}

	searchViewOffset := hiddenViewOffset
	if gui.State.Searching.isSearching {
		searchViewOffset = 0
	}

	// this view takes up one character. Its only purpose is to show the slash when searching
	searchPrefix := "search: "
	if searchPrefixView, err := g.SetView("searchPrefix", appStatusOptionsBoundary-1+searchViewOffset, height-2+searchViewOffset, len(searchPrefix)+searchViewOffset, height+searchViewOffset, 0); err != nil {
		if err.Error() != "unknown view" {
			return err
		}

		searchPrefixView.BgColor = gocui.ColorDefault
		searchPrefixView.FgColor = gocui.ColorGreen
		searchPrefixView.Frame = false
		gui.setViewContent(searchPrefixView, searchPrefix)
	}

	if searchView, err := g.SetView("search", appStatusOptionsBoundary-1+searchViewOffset+len(searchPrefix), height-2+searchViewOffset, optionsVersionBoundary+searchViewOffset, height+searchViewOffset, 0); err != nil {
		if err.Error() != "unknown view" {
			return err
		}

		searchView.BgColor = gocui.ColorDefault
		searchView.FgColor = gocui.ColorGreen
		searchView.Frame = false
		searchView.Editable = true
	}

	if appStatusView, err := g.SetView("appStatus", -1, height-2, width, height, 0); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		appStatusView.BgColor = gocui.ColorDefault
		appStatusView.FgColor = gocui.ColorCyan
		appStatusView.Frame = false
		if _, err := g.SetViewOnBottom("appStatus"); err != nil {
			return err
		}
	}

	informationView, err := g.SetView("information", optionsVersionBoundary-1, height-2, width, height, 0)
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		informationView.BgColor = gocui.ColorDefault
		informationView.FgColor = gocui.ColorGreen
		informationView.Frame = false
		gui.renderString(g, "information", information)

		// doing this here because it'll only happen once
		if err := gui.onInitialViewsCreation(); err != nil {
			return err
		}
	}
	if gui.State.OldInformation != information {
		gui.setViewContent(informationView, information)
		gui.State.OldInformation = information
	}

	if gui.g.CurrentView() == nil {
		initialView := gui.getFilesView()
		if gui.inFilterMode() {
			initialView = gui.getCommitsView()
		}
		if _, err := gui.g.SetCurrentView(initialView.Name()); err != nil {
			return err
		}

		if err := gui.switchFocus(gui.g, nil, initialView); err != nil {
			return err
		}
	}

	type listViewState struct {
		selectedLine int
		lineCount    int
		view         *gocui.View
		context      string
	}

	listViews := []listViewState{
		{view: filesView, context: "", selectedLine: gui.State.Panels.Files.SelectedLine, lineCount: len(gui.State.Files)},
		{view: branchesView, context: "local-branches", selectedLine: gui.State.Panels.Branches.SelectedLine, lineCount: len(gui.State.Branches)},
		{view: branchesView, context: "remotes", selectedLine: gui.State.Panels.Remotes.SelectedLine, lineCount: len(gui.State.Remotes)},
		{view: branchesView, context: "remote-branches", selectedLine: gui.State.Panels.RemoteBranches.SelectedLine, lineCount: len(gui.State.Remotes)},
		{view: commitsView, context: "branch-commits", selectedLine: gui.State.Panels.Commits.SelectedLine, lineCount: len(gui.State.Commits)},
		{view: commitsView, context: "reflog-commits", selectedLine: gui.State.Panels.ReflogCommits.SelectedLine, lineCount: len(gui.State.FilteredReflogCommits)},
		{view: stashView, context: "", selectedLine: gui.State.Panels.Stash.SelectedLine, lineCount: len(gui.State.StashEntries)},
		{view: commitFilesView, context: "", selectedLine: gui.State.Panels.CommitFiles.SelectedLine, lineCount: len(gui.State.CommitFiles)},
	}

	// menu view might not exist so we check to be safe
	if menuView, err := gui.g.View("menu"); err == nil {
		listViews = append(listViews, listViewState{view: menuView, context: "", selectedLine: gui.State.Panels.Menu.SelectedLine, lineCount: gui.State.MenuItemCount})
	}
	for _, listView := range listViews {
		// ignore views where the context doesn't match up with the selected line we're trying to focus
		if listView.context != "" && (listView.view.Context != listView.context) {
			continue
		}
		// check if the selected line is now out of view and if so refocus it
		listView.view.FocusPoint(0, listView.selectedLine)

		listView.view.SelBgColor = theme.GocuiSelectedLineBgColor
	}

	mainViewWidth, mainViewHeight := gui.getMainView().Size()
	if mainViewWidth != gui.State.PrevMainWidth || mainViewHeight != gui.State.PrevMainHeight {
		gui.State.PrevMainWidth = mainViewWidth
		gui.State.PrevMainHeight = mainViewHeight
		if err := gui.onResize(); err != nil {
			return err
		}
	}

	// here is a good place log some stuff
	// if you download humanlog and do tail -f development.log | humanlog
	// this will let you see these branches as prettified json
	// gui.Log.Info(utils.AsJson(gui.State.Branches[0:4]))
	return gui.resizeCurrentPopupPanel(g)
}

func (gui *Gui) onInitialViewsCreation() error {
	gui.changeMainViewsContext("normal")

	gui.getBranchesView().Context = "local-branches"
	gui.getCommitsView().Context = "branch-commits"

	return gui.loadNewRepo()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
