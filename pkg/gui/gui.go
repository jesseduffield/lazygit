package gui

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"runtime"
	"sync"

	// "io"
	// "io/ioutil"

	"os/exec"
	"strings"
	"time"

	"github.com/go-errors/errors"

	// "strings"

	"github.com/fatih/color"
	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/mattn/go-runewidth"
	"github.com/sirupsen/logrus"
)

const (
	SCREEN_NORMAL int = iota
	SCREEN_HALF
	SCREEN_FULL
)

const StartupPopupVersion = 1

// OverlappingEdges determines if panel edges overlap
var OverlappingEdges = false

// SentinelErrors are the errors that have special meaning and need to be checked
// by calling functions. The less of these, the better
type SentinelErrors struct {
	ErrSubProcess error
	ErrNoFiles    error
	ErrSwitchRepo error
}

// GenerateSentinelErrors makes the sentinel errors for the gui. We're defining it here
// because we can't do package-scoped errors with localization, and also because
// it seems like package-scoped variables are bad in general
// https://dave.cheney.net/2017/06/11/go-without-package-scoped-variables
// In the future it would be good to implement some of the recommendations of
// that article. For now, if we don't need an error to be a sentinel, we will just
// define it inline. This has implications for error messages that pop up everywhere
// in that we'll be duplicating the default values. We may need to look at
// having a default localisation bundle defined, and just using keys-only when
// localising things in the code.
func (gui *Gui) GenerateSentinelErrors() {
	gui.Errors = SentinelErrors{
		ErrSubProcess: errors.New(gui.Tr.SLocalize("RunningSubprocess")),
		ErrNoFiles:    errors.New(gui.Tr.SLocalize("NoChangedFiles")),
		ErrSwitchRepo: errors.New("switching repo"),
	}
}

// Teml is short for template used to make the required map[string]interface{} shorter when using gui.Tr.SLocalize and gui.Tr.TemplateLocalize
type Teml i18n.Teml

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	g                    *gocui.Gui
	Log                  *logrus.Entry
	GitCommand           *commands.GitCommand
	OSCommand            *commands.OSCommand
	SubProcess           *exec.Cmd
	State                guiState
	Config               config.AppConfigurer
	Tr                   *i18n.Localizer
	Errors               SentinelErrors
	Updater              *updates.Updater
	statusManager        *statusManager
	credentials          credentials
	waitForIntro         sync.WaitGroup
	fileWatcher          *fileWatcher
	viewBufferManagerMap map[string]*tasks.ViewBufferManager
	stopChan             chan struct{}
}

// for now the staging panel state, unlike the other panel states, is going to be
// non-mutative, so that we don't accidentally end up
// with mismatches of data. We might change this in the future
type lineByLinePanelState struct {
	SelectedLineIdx  int
	FirstLineIdx     int
	LastLineIdx      int
	Diff             string
	PatchParser      *commands.PatchParser
	SelectMode       int  // one of LINE, HUNK, or RANGE
	SecondaryFocused bool // this is for if we show the left or right panel
}

type mergingPanelState struct {
	ConflictIndex int
	ConflictTop   bool
	Conflicts     []commands.Conflict
	EditHistory   *stack.Stack
}

type filePanelState struct {
	SelectedLine int
}

// TODO: consider splitting this out into the window and the branches view
type branchPanelState struct {
	SelectedLine int
}

type remotePanelState struct {
	SelectedLine int
}

type remoteBranchesState struct {
	SelectedLine int
}

type tagsPanelState struct {
	SelectedLine int
}

type commitPanelState struct {
	SelectedLine     int
	SpecificDiffMode bool
	LimitCommits     bool
}

type reflogCommitPanelState struct {
	SelectedLine int
}

type stashPanelState struct {
	SelectedLine int
}

type menuPanelState struct {
	SelectedLine int
	OnPress      func(g *gocui.Gui, v *gocui.View) error
}

type commitFilesPanelState struct {
	SelectedLine int
}

type statusPanelState struct {
	pushables string
	pullables string
}

type panelStates struct {
	Files          *filePanelState
	Branches       *branchPanelState
	Remotes        *remotePanelState
	RemoteBranches *remoteBranchesState
	Tags           *tagsPanelState
	Commits        *commitPanelState
	ReflogCommits  *reflogCommitPanelState
	Stash          *stashPanelState
	Menu           *menuPanelState
	LineByLine     *lineByLinePanelState
	Merging        *mergingPanelState
	CommitFiles    *commitFilesPanelState
	Status         *statusPanelState
}

type searchingState struct {
	view         *gocui.View
	isSearching  bool
	searchString string
}

type guiState struct {
	Files                []*commands.File
	Branches             []*commands.Branch
	Commits              []*commands.Commit
	StashEntries         []*commands.StashEntry
	CommitFiles          []*commands.CommitFile
	ReflogCommits        []*commands.Commit
	DiffEntries          []*commands.Commit
	Remotes              []*commands.Remote
	RemoteBranches       []*commands.RemoteBranch
	Tags                 []*commands.Tag
	MenuItemCount        int // can't store the actual list because it's of interface{} type
	PreviousView         string
	Platform             commands.Platform
	Updating             bool
	Panels               *panelStates
	WorkingTreeState     string // one of "merging", "rebasing", "normal"
	MainContext          string // used to keep the main and secondary views' contexts in sync
	CherryPickedCommits  []*commands.Commit
	SplitMainPanel       bool
	RetainOriginalDir    bool
	IsRefreshingFiles    bool
	RefreshingFilesMutex sync.Mutex
	Searching            searchingState
	ScreenMode           int
	SideView             *gocui.View
}

// for now the split view will always be on

// NewGui builds a new gui handler
func NewGui(log *logrus.Entry, gitCommand *commands.GitCommand, oSCommand *commands.OSCommand, tr *i18n.Localizer, config config.AppConfigurer, updater *updates.Updater) (*Gui, error) {

	initialState := guiState{
		Files:               make([]*commands.File, 0),
		PreviousView:        "files",
		Commits:             make([]*commands.Commit, 0),
		CherryPickedCommits: make([]*commands.Commit, 0),
		StashEntries:        make([]*commands.StashEntry, 0),
		DiffEntries:         make([]*commands.Commit, 0),
		Platform:            *oSCommand.Platform,
		Panels: &panelStates{
			Files:          &filePanelState{SelectedLine: -1},
			Branches:       &branchPanelState{SelectedLine: 0},
			Remotes:        &remotePanelState{SelectedLine: 0},
			RemoteBranches: &remoteBranchesState{SelectedLine: -1},
			Tags:           &tagsPanelState{SelectedLine: -1},
			Commits:        &commitPanelState{SelectedLine: -1, LimitCommits: true},
			ReflogCommits:  &reflogCommitPanelState{SelectedLine: 0}, // TODO: might need to make -1
			CommitFiles:    &commitFilesPanelState{SelectedLine: -1},
			Stash:          &stashPanelState{SelectedLine: -1},
			Menu:           &menuPanelState{SelectedLine: 0},
			Merging: &mergingPanelState{
				ConflictIndex: 0,
				ConflictTop:   true,
				Conflicts:     []commands.Conflict{},
				EditHistory:   stack.New(),
			},
			Status: &statusPanelState{},
		},
		ScreenMode: SCREEN_NORMAL,
		SideView:   nil,
	}

	gui := &Gui{
		Log:                  log,
		GitCommand:           gitCommand,
		OSCommand:            oSCommand,
		State:                initialState,
		Config:               config,
		Tr:                   tr,
		Updater:              updater,
		statusManager:        &statusManager{},
		viewBufferManagerMap: map[string]*tasks.ViewBufferManager{},
	}

	gui.watchFilesForChanges()

	gui.GenerateSentinelErrors()

	return gui, nil
}

func (gui *Gui) nextScreenMode(g *gocui.Gui, v *gocui.View) error {
	gui.State.ScreenMode = utils.NextIntInCycle([]int{SCREEN_NORMAL, SCREEN_HALF, SCREEN_FULL}, gui.State.ScreenMode)
	// commits render differently depending on whether we're in fullscreen more or not
	if err := gui.refreshCommitsViewWithSelection(); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) prevScreenMode(g *gocui.Gui, v *gocui.View) error {
	gui.State.ScreenMode = utils.PrevIntInCycle([]int{SCREEN_NORMAL, SCREEN_HALF, SCREEN_FULL}, gui.State.ScreenMode)
	// commits render differently depending on whether we're in fullscreen more or not
	if err := gui.refreshCommitsViewWithSelection(); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) scrollUpView(viewName string) error {
	mainView, _ := gui.g.View(viewName)
	ox, oy := mainView.Origin()
	newOy := int(math.Max(0, float64(oy-gui.Config.GetUserConfig().GetInt("gui.scrollHeight"))))
	return mainView.SetOrigin(ox, newOy)
}

func (gui *Gui) scrollDownView(viewName string) error {
	mainView, _ := gui.g.View(viewName)
	ox, oy := mainView.Origin()
	y := oy
	if !gui.Config.GetUserConfig().GetBool("gui.scrollPastBottom") {
		_, sy := mainView.Size()
		y += sy
	}
	scrollHeight := gui.Config.GetUserConfig().GetInt("gui.scrollHeight")
	if y < mainView.LinesHeight() {
		if err := mainView.SetOrigin(ox, oy+scrollHeight); err != nil {
			return err
		}
	}
	if manager, ok := gui.viewBufferManagerMap[viewName]; ok {
		manager.ReadLines(scrollHeight)
	}
	return nil
}

func (gui *Gui) scrollUpMain(g *gocui.Gui, v *gocui.View) error {
	return gui.scrollUpView("main")
}

func (gui *Gui) scrollDownMain(g *gocui.Gui, v *gocui.View) error {
	return gui.scrollDownView("main")
}

func (gui *Gui) scrollUpSecondary(g *gocui.Gui, v *gocui.View) error {
	return gui.scrollUpView("secondary")
}

func (gui *Gui) scrollDownSecondary(g *gocui.Gui, v *gocui.View) error {
	return gui.scrollDownView("secondary")
}

func (gui *Gui) handleRefresh(g *gocui.Gui, v *gocui.View) error {
	return gui.refreshSidePanels(g)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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
		gui.onSearchEscape()
	}
	switch v.Name() {
	case "branches":
		if v.Context == "local-branches" {
			// This stops the branches panel from showing the upstream/downstream changes to the selected branch, when it loses focus
			displayStrings := presentation.GetBranchListDisplayStrings(gui.State.Branches, false, -1)
			gui.renderDisplayStrings(gui.getBranchesView(), displayStrings)
		}
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
	g.DeleteView("limit")

	textColor := theme.GocuiDefaultTextColor
	var leftSideWidth int
	switch gui.State.ScreenMode {
	case SCREEN_NORMAL:
		leftSideWidth = width / 3
	case SCREEN_HALF:
		leftSideWidth = width / 2
	case SCREEN_FULL:
		currentView := gui.g.CurrentView()
		if currentView != nil && currentView.Name() == "main" {
			leftSideWidth = 0
		} else {
			leftSideWidth = width - 1
		}
	}

	panelSplitX := width - 1
	mainPanelLeft := leftSideWidth + 1
	mainPanelRight := width - 1
	secondaryPanelLeft := width - 1
	secondaryPanelTop := 0
	mainPanelBottom := height - 2
	if gui.State.SplitMainPanel {
		if gui.State.ScreenMode == SCREEN_FULL {
			mainPanelLeft = 0
			panelSplitX = width/2 - 4
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
			panelSplitX = (1 + ((units - 1) / 2)) * width / units
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

	if v, err := g.SetViewBeneath("commitFiles", "branches", vHeights["commits"]); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.Title = gui.Tr.SLocalize("CommitFiles")
		v.FgColor = textColor
		v.SetOnSelectItem(gui.onSelectItemWrapper(gui.onCommitFilesPanelSearchSelect))
		v.ContainsList = true
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
			g.SetViewOnBottom("commitMessage")
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
		gui.setViewContent(gui.g, searchPrefixView, searchPrefix)
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

	if v, err := g.SetView("information", optionsVersionBoundary-1, height-2, width, height, 0); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.BgColor = gocui.ColorDefault
		v.FgColor = gocui.ColorGreen
		v.Frame = false
		if err := gui.renderString(g, "information", information); err != nil {
			return err
		}

		// doing this here because it'll only happen once
		if err := gui.onInitialViewsCreation(); err != nil {
			return err
		}
	}

	if gui.g.CurrentView() == nil {
		if _, err := gui.g.SetCurrentView(gui.getFilesView().Name()); err != nil {
			return err
		}

		if err := gui.switchFocus(gui.g, nil, gui.getFilesView()); err != nil {
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
		{view: commitsView, context: "reflog-commits", selectedLine: gui.State.Panels.ReflogCommits.SelectedLine, lineCount: len(gui.State.ReflogCommits)},
		{view: stashView, context: "", selectedLine: gui.State.Panels.Stash.SelectedLine, lineCount: len(gui.State.StashEntries)},
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
		if err := gui.focusPoint(0, listView.selectedLine, listView.lineCount, listView.view); err != nil {
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

func (gui *Gui) loadNewRepo() error {
	gui.Updater.CheckForNewUpdate(gui.onBackgroundUpdateCheckFinish, false)
	if err := gui.updateRecentRepoList(); err != nil {
		return err
	}
	gui.waitForIntro.Done()

	if err := gui.refreshSidePanels(gui.g); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) showInitialPopups(tasks []func(chan struct{}) error) {
	gui.waitForIntro.Add(len(tasks))
	done := make(chan struct{})

	go func() {
		for _, task := range tasks {
			go func() {
				if err := task(done); err != nil {
					_ = gui.createErrorPanel(gui.g, err.Error())
				}
			}()

			<-done
			gui.waitForIntro.Done()
		}
	}()
}

func (gui *Gui) showShamelessSelfPromotionMessage(done chan struct{}) error {
	onConfirm := func(g *gocui.Gui, v *gocui.View) error {
		done <- struct{}{}
		return gui.Config.WriteToUserConfig("startupPopupVersion", StartupPopupVersion)
	}

	return gui.createConfirmationPanel(gui.g, nil, true, gui.Tr.SLocalize("ShamelessSelfPromotionTitle"), gui.Tr.SLocalize("ShamelessSelfPromotionMessage"), onConfirm, onConfirm)
}

func (gui *Gui) promptAnonymousReporting(done chan struct{}) error {
	return gui.createConfirmationPanel(gui.g, nil, true, gui.Tr.SLocalize("AnonymousReportingTitle"), gui.Tr.SLocalize("AnonymousReportingPrompt"), func(g *gocui.Gui, v *gocui.View) error {
		done <- struct{}{}
		return gui.Config.WriteToUserConfig("reporting", "on")
	}, func(g *gocui.Gui, v *gocui.View) error {
		done <- struct{}{}
		return gui.Config.WriteToUserConfig("reporting", "off")
	})
}

func (gui *Gui) fetch(g *gocui.Gui, v *gocui.View, canAskForCredentials bool) (unamePassOpend bool, err error) {
	unamePassOpend = false
	err = gui.GitCommand.Fetch(func(passOrUname string) string {
		unamePassOpend = true
		return gui.waitForPassUname(gui.g, v, passOrUname)
	}, canAskForCredentials)

	if canAskForCredentials && err != nil && strings.Contains(err.Error(), "exit status 128") {
		colorFunction := color.New(color.FgRed).SprintFunc()
		coloredMessage := colorFunction(strings.TrimSpace(gui.Tr.SLocalize("PassUnameWrong")))
		close := func(g *gocui.Gui, v *gocui.View) error {
			return nil
		}
		_ = gui.createConfirmationPanel(g, v, true, gui.Tr.SLocalize("Error"), coloredMessage, close, close)
	}

	gui.refreshStatus(g)
	return unamePassOpend, err
}

func (gui *Gui) renderGlobalOptions() error {
	return gui.renderOptionsMap(map[string]string{
		fmt.Sprintf("%s/%s", gui.getKeyDisplay("universal.scrollUpMain"), gui.getKeyDisplay("universal.scrollDownMain")):                                                                                 gui.Tr.SLocalize("scroll"),
		fmt.Sprintf("%s %s %s %s", gui.getKeyDisplay("universal.prevBlock"), gui.getKeyDisplay("universal.nextBlock"), gui.getKeyDisplay("universal.prevItem"), gui.getKeyDisplay("universal.nextItem")): gui.Tr.SLocalize("navigate"),
		fmt.Sprintf("%s/%s", gui.getKeyDisplay("universal.return"), gui.getKeyDisplay("universal.quit")):                                                                                                 gui.Tr.SLocalize("close"),
		fmt.Sprintf("%s", gui.getKeyDisplay("universal.optionMenu")):                                                                                                                                     gui.Tr.SLocalize("menu"),
		"1-5": gui.Tr.SLocalize("jump"),
	})
}

func (gui *Gui) goEvery(interval time.Duration, stop chan struct{}, function func() error) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = function()
			case <-stop:
				return
			}
		}
	}()
}

func (gui *Gui) startBackgroundFetch() {
	gui.waitForIntro.Wait()
	isNew := gui.Config.GetIsNewRepo()
	if !isNew {
		time.After(60 * time.Second)
	}
	_, err := gui.fetch(gui.g, gui.g.CurrentView(), false)
	if err != nil && strings.Contains(err.Error(), "exit status 128") && isNew {
		_ = gui.createConfirmationPanel(gui.g, gui.g.CurrentView(), true, gui.Tr.SLocalize("NoAutomaticGitFetchTitle"), gui.Tr.SLocalize("NoAutomaticGitFetchBody"), nil, nil)
	} else {
		gui.goEvery(time.Second*60, gui.stopChan, func() error {
			_, err := gui.fetch(gui.g, gui.g.CurrentView(), false)
			return err
		})
	}
}

// Run setup the gui with keybindings and start the mainloop
func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal, OverlappingEdges)
	if err != nil {
		return err
	}
	defer g.Close()

	g.OnSearchEscape = gui.onSearchEscape
	g.SearchEscapeKey = gui.getKey("universal.return")
	g.NextSearchMatchKey = gui.getKey("universal.nextMatch")
	g.PrevSearchMatchKey = gui.getKey("universal.prevMatch")

	gui.stopChan = make(chan struct{})

	g.ASCII = runtime.GOOS == "windows" && runewidth.IsEastAsian()

	if gui.Config.GetUserConfig().GetBool("gui.mouseEvents") {
		g.Mouse = true
	}

	gui.g = g // TODO: always use gui.g rather than passing g around everywhere

	if err := gui.setColorScheme(); err != nil {
		return err
	}

	popupTasks := []func(chan struct{}) error{}
	if gui.Config.GetUserConfig().GetString("reporting") == "undetermined" {
		popupTasks = append(popupTasks, gui.promptAnonymousReporting)
	}
	configPopupVersion := gui.Config.GetUserConfig().GetInt("StartupPopupVersion")
	// -1 means we've disabled these popups
	if configPopupVersion != -1 && configPopupVersion < StartupPopupVersion {
		popupTasks = append(popupTasks, gui.showShamelessSelfPromotionMessage)
	}
	gui.showInitialPopups(popupTasks)

	gui.waitForIntro.Add(1)
	if gui.Config.GetUserConfig().GetBool("git.autoFetch") {
		go gui.startBackgroundFetch()
	}

	gui.goEvery(time.Second*10, gui.stopChan, gui.refreshFiles)

	g.SetManager(gocui.ManagerFunc(gui.layout), gocui.ManagerFunc(gui.getFocusLayout()))

	if err = gui.keybindings(g); err != nil {
		return err
	}

	gui.Log.Warn("starting main loop")

	err = g.MainLoop()
	return err
}

// RunWithSubprocesses loops, instantiating a new gocui.Gui with each iteration
// if the error returned from a run is a ErrSubProcess, it runs the subprocess
// otherwise it handles the error, possibly by quitting the application
func (gui *Gui) RunWithSubprocesses() error {
	for {
		if err := gui.Run(); err != nil {
			for _, manager := range gui.viewBufferManagerMap {
				manager.Close()
			}
			gui.viewBufferManagerMap = map[string]*tasks.ViewBufferManager{}

			if !gui.fileWatcher.Disabled {
				gui.fileWatcher.Watcher.Close()
			}

			close(gui.stopChan)

			if err == gocui.ErrQuit {
				if !gui.State.RetainOriginalDir {
					if err := gui.recordCurrentDirectory(); err != nil {
						return err
					}
				}

				break
			} else if err == gui.Errors.ErrSwitchRepo {
				continue
			} else if err == gui.Errors.ErrSubProcess {
				if err := gui.runCommand(); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func (gui *Gui) runCommand() error {
	gui.SubProcess.Stdout = os.Stdout
	gui.SubProcess.Stderr = os.Stdout
	gui.SubProcess.Stdin = os.Stdin

	fmt.Fprintf(os.Stdout, "\n%s\n\n", utils.ColoredString("+ "+strings.Join(gui.SubProcess.Args, " "), color.FgBlue))

	if err := gui.SubProcess.Run(); err != nil {
		// not handling the error explicitly because usually we're going to see it
		// in the output anyway
		gui.Log.Error(err)
	}

	gui.SubProcess.Stdout = ioutil.Discard
	gui.SubProcess.Stderr = ioutil.Discard
	gui.SubProcess.Stdin = nil
	gui.SubProcess = nil

	fmt.Fprintf(os.Stdout, "\n%s", utils.ColoredString(gui.Tr.SLocalize("pressEnterToReturn"), color.FgGreen))
	fmt.Scanln() // wait for enter press

	return nil
}

func (gui *Gui) handleDonate(g *gocui.Gui, v *gocui.View) error {
	if !gui.g.Mouse {
		return nil
	}

	cx, _ := v.Cursor()
	if cx > len(gui.Tr.SLocalize("Donate")) {
		return nil
	}
	return gui.OSCommand.OpenLink("https://github.com/sponsors/jesseduffield")
}

// setColorScheme sets the color scheme for the app based on the user config
func (gui *Gui) setColorScheme() error {
	userConfig := gui.Config.GetUserConfig()
	theme.UpdateTheme(userConfig)

	gui.g.FgColor = theme.InactiveBorderColor
	gui.g.SelFgColor = theme.ActiveBorderColor

	return nil
}

func (gui *Gui) handleMouseDownMain(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	switch g.CurrentView().Name() {
	case "files":
		return gui.enterFile(false, v.SelectedLineIdx())
	case "commitFiles":
		return gui.enterCommitFile(v.SelectedLineIdx())
	}

	return nil
}

func (gui *Gui) handleMouseDownSecondary(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	switch g.CurrentView().Name() {
	case "files":
		return gui.enterFile(true, v.SelectedLineIdx())
	}

	return nil
}
