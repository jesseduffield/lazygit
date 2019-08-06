package gui

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
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
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

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
	g             *gocui.Gui
	Log           *logrus.Entry
	GitCommand    *commands.GitCommand
	OSCommand     *commands.OSCommand
	SubProcess    *exec.Cmd
	State         guiState
	Config        config.AppConfigurer
	Tr            *i18n.Localizer
	Errors        SentinelErrors
	Updater       *updates.Updater
	statusManager *statusManager
	credentials   credentials
	waitForIntro  sync.WaitGroup
}

// for now the staging panel state, unlike the other panel states, is going to be
// non-mutative, so that we don't accidentally end up
// with mismatches of data. We might change this in the future
type stagingPanelState struct {
	SelectedLine   int
	StageableLines []int
	HunkStarts     []int
	Diff           string
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

type branchPanelState struct {
	SelectedLine int
}

type commitPanelState struct {
	SelectedLine     int
	SpecificDiffMode bool
}

type stashPanelState struct {
	SelectedLine int
}

type menuPanelState struct {
	SelectedLine int
}

type commitFilesPanelState struct {
	SelectedLine int
}

type panelStates struct {
	Files       *filePanelState
	Branches    *branchPanelState
	Commits     *commitPanelState
	Stash       *stashPanelState
	Menu        *menuPanelState
	Staging     *stagingPanelState
	Merging     *mergingPanelState
	CommitFiles *commitFilesPanelState
}

type guiState struct {
	Files               []*commands.File
	Branches            []*commands.Branch
	Commits             []*commands.Commit
	StashEntries        []*commands.StashEntry
	CommitFiles         []*commands.CommitFile
	DiffEntries         []*commands.Commit
	MenuItemCount       int // can't store the actual list because it's of interface{} type
	PreviousView        string
	Platform            commands.Platform
	Updating            bool
	Panels              *panelStates
	WorkingTreeState    string // one of "merging", "rebasing", "normal"
	Contexts            map[string]string
	CherryPickedCommits []*commands.Commit
}

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
			Files:       &filePanelState{SelectedLine: -1},
			Branches:    &branchPanelState{SelectedLine: 0},
			Commits:     &commitPanelState{SelectedLine: -1},
			CommitFiles: &commitFilesPanelState{SelectedLine: -1},
			Stash:       &stashPanelState{SelectedLine: -1},
			Menu:        &menuPanelState{SelectedLine: 0},
			Merging: &mergingPanelState{
				ConflictIndex: 0,
				ConflictTop:   true,
				Conflicts:     []commands.Conflict{},
				EditHistory:   stack.New(),
			},
		},
	}

	gui := &Gui{
		Log:           log,
		GitCommand:    gitCommand,
		OSCommand:     oSCommand,
		State:         initialState,
		Config:        config,
		Tr:            tr,
		Updater:       updater,
		statusManager: &statusManager{},
	}

	gui.GenerateSentinelErrors()

	return gui, nil
}

func (gui *Gui) scrollUpMain(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := g.View("main")
	ox, oy := mainView.Origin()
	newOy := int(math.Max(0, float64(oy-gui.Config.GetUserConfig().GetInt("gui.scrollHeight"))))
	return mainView.SetOrigin(ox, newOy)
}

func (gui *Gui) scrollDownMain(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := g.View("main")
	ox, oy := mainView.Origin()
	y := oy
	if !gui.Config.GetUserConfig().GetBool("gui.scrollPastBottom") {
		_, sy := mainView.Size()
		y += sy
	}
	if y < len(mainView.BufferLines()) {
		return mainView.SetOrigin(ox, oy+gui.Config.GetUserConfig().GetInt("gui.scrollHeight"))
	}
	return nil
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
	return gui.setMainTitle()
}

func (gui *Gui) onFocusLost(v *gocui.View, newView *gocui.View) error {
	if v == nil {
		return nil
	}
	if v.Name() == "branches" {
		// This stops the branches panel from showing the upstream/downstream changes to the selected branch, when it loses focus
		// inside renderListPanel it checks to see if the panel has focus
		if err := gui.renderListPanel(gui.getBranchesView(), gui.State.Branches); err != nil {
			return err
		}
	} else if v.Name() == "main" {
		// if we have lost focus to a first-class panel, we need to do some cleanup
		if err := gui.changeContext("main", "normal"); err != nil {
			return err
		}

	} else if v.Name() == "commitFiles" {
		if _, err := gui.g.SetViewOnBottom(v.Name()); err != nil {
			return err
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

	currView := gui.g.CurrentView()
	currentCyclebleView := gui.State.PreviousView
	if currView != nil {
		viewName := currView.Name()
		usePreviouseView := true
		for _, view := range cyclableViews {
			if view == viewName {
				currentCyclebleView = viewName
				usePreviouseView = false
				break
			}
		}
		if usePreviouseView {
			currentCyclebleView = gui.State.PreviousView
		}
	}

	usableSpace := height - 7
	extraSpace := usableSpace - (usableSpace/3)*3

	vHeights := map[string]int{
		"status":   3,
		"files":    (usableSpace / 3) + extraSpace,
		"branches": usableSpace / 3,
		"commits":  usableSpace / 3,
		"stash":    3,
		"options":  1,
	}

	if height < 28 {
		defaultHeight := 3
		if height < 21 {
			defaultHeight = 1
		}
		vHeights = map[string]int{
			"status":   defaultHeight,
			"files":    defaultHeight,
			"branches": defaultHeight,
			"commits":  defaultHeight,
			"stash":    defaultHeight,
			"options":  defaultHeight,
		}
		vHeights[currentCyclebleView] = height - defaultHeight*4 - 1
	}

	optionsVersionBoundary := width - max(len(utils.Decolorise(information)), 1)
	leftSideWidth := width / 3

	appStatus := gui.statusManager.getStatusString()
	appStatusOptionsBoundary := 0
	if appStatus != "" {
		appStatusOptionsBoundary = len(appStatus) + 2
	}

	panelSpacing := 1
	if OverlappingEdges {
		panelSpacing = 0
	}

	_, _ = g.SetViewOnBottom("limit")
	g.DeleteView("limit")

	v, err := g.SetView("main", leftSideWidth+panelSpacing, 0, width-1, height-2, gocui.LEFT)
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.Title = gui.Tr.SLocalize("DiffTitle")
		v.Wrap = true
		v.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("status", 0, 0, leftSideWidth, vHeights["status"]-1, gocui.BOTTOM|gocui.RIGHT); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.Title = gui.Tr.SLocalize("StatusTitle")
		v.FgColor = gocui.ColorWhite
	}

	filesView, err := g.SetViewBeneath("files", "status", vHeights["files"])
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		filesView.Highlight = true
		filesView.Title = gui.Tr.SLocalize("FilesTitle")
		v.FgColor = gocui.ColorWhite
	}

	branchesView, err := g.SetViewBeneath("branches", "files", vHeights["branches"])
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		branchesView.Title = gui.Tr.SLocalize("BranchesTitle")
		branchesView.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetViewBeneath("commitFiles", "branches", vHeights["commits"]); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.Title = gui.Tr.SLocalize("CommitFiles")
		v.FgColor = gocui.ColorWhite
	}

	commitsView, err := g.SetViewBeneath("commits", "branches", vHeights["commits"])
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		commitsView.Title = gui.Tr.SLocalize("CommitsTitle")
		commitsView.FgColor = gocui.ColorWhite
	}

	stashView, err := g.SetViewBeneath("stash", "commits", vHeights["stash"])
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		stashView.Title = gui.Tr.SLocalize("StashTitle")
		stashView.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("options", appStatusOptionsBoundary-1, height-2, optionsVersionBoundary-1, height, 0); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.Frame = false
		if v.FgColor, err = gui.GetOptionsPanelTextColor(); err != nil {
			return err
		}
	}

	if gui.getCommitMessageView() == nil {
		// doesn't matter where this view starts because it will be hidden
		if commitMessageView, err := g.SetView("commitMessage", width, height, width*2, height*2, 0); err != nil {
			if err.Error() != "unknown view" {
				return err
			}
			g.SetViewOnBottom("commitMessage")
			commitMessageView.Title = gui.Tr.SLocalize("CommitMessage")
			commitMessageView.FgColor = gocui.ColorWhite
			commitMessageView.Editable = true
		}
	}

	if check, _ := g.View("credentials"); check == nil {
		// doesn't matter where this view starts because it will be hidden
		if credentialsView, err := g.SetView("credentials", width, height, width*2, height*2, 0); err != nil {
			if err.Error() != "unknown view" {
				return err
			}
			_, err := g.SetViewOnBottom("credentials")
			if err != nil {
				return err
			}
			credentialsView.Title = gui.Tr.SLocalize("CredentialsUsername")
			credentialsView.FgColor = gocui.ColorWhite
			credentialsView.Editable = true
		}
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
		if err := gui.loadNewRepo(); err != nil {
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
	}

	listViews := map[*gocui.View]listViewState{
		filesView:    {selectedLine: gui.State.Panels.Files.SelectedLine, lineCount: len(gui.State.Files)},
		branchesView: {selectedLine: gui.State.Panels.Branches.SelectedLine, lineCount: len(gui.State.Branches)},
		commitsView:  {selectedLine: gui.State.Panels.Commits.SelectedLine, lineCount: len(gui.State.Commits)},
		stashView:    {selectedLine: gui.State.Panels.Stash.SelectedLine, lineCount: len(gui.State.StashEntries)},
	}

	// menu view might not exist so we check to be safe
	if menuView, err := gui.g.View("menu"); err == nil {
		listViews[menuView] = listViewState{selectedLine: gui.State.Panels.Menu.SelectedLine, lineCount: gui.State.MenuItemCount}
	}
	for view, state := range listViews {
		// check if the selected line is now out of view and if so refocus it
		if err := gui.focusPoint(0, state.selectedLine, state.lineCount, view); err != nil {
			return err
		}
	}

	// here is a good place log some stuff
	// if you download humanlog and do tail -f development.log | humanlog
	// this will let you see these branches as prettified json
	// gui.Log.Info(utils.AsJson(gui.State.Branches[0:4]))
	return gui.resizeCurrentPopupPanel(g)
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

	if gui.Config.GetUserConfig().GetString("reporting") == "undetermined" {
		if err := gui.promptAnonymousReporting(); err != nil {
			return err
		}
	}
	return nil
}

func (gui *Gui) promptAnonymousReporting() error {
	return gui.createConfirmationPanel(gui.g, nil, gui.Tr.SLocalize("AnonymousReportingTitle"), gui.Tr.SLocalize("AnonymousReportingPrompt"), func(g *gocui.Gui, v *gocui.View) error {
		gui.waitForIntro.Done()
		return gui.Config.WriteToUserConfig("reporting", "on")
	}, func(g *gocui.Gui, v *gocui.View) error {
		gui.waitForIntro.Done()
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
		_ = gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("Error"), coloredMessage, close, close)
	}

	gui.refreshStatus(g)
	return unamePassOpend, err
}

func (gui *Gui) renderAppStatus() error {
	appStatus := gui.statusManager.getStatusString()
	if appStatus != "" {
		return gui.renderString(gui.g, "appStatus", appStatus)
	}
	return nil
}

func (gui *Gui) renderGlobalOptions() error {
	return gui.renderOptionsMap(map[string]string{
		"PgUp/PgDn": gui.Tr.SLocalize("scroll"),
		"← → ↑ ↓":   gui.Tr.SLocalize("navigate"),
		"esc/q":     gui.Tr.SLocalize("close"),
		"x":         gui.Tr.SLocalize("menu"),
	})
}

func (gui *Gui) goEvery(interval time.Duration, function func() error) {
	go func() {
		for range time.Tick(interval) {
			_ = function()
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
		_ = gui.createConfirmationPanel(gui.g, gui.g.CurrentView(), gui.Tr.SLocalize("NoAutomaticGitFetchTitle"), gui.Tr.SLocalize("NoAutomaticGitFetchBody"), nil, nil)
	} else {
		gui.goEvery(time.Second*60, func() error {
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

	if gui.Config.GetUserConfig().GetBool("gui.mouseEvents") {
		g.Mouse = true
	}

	gui.g = g // TODO: always use gui.g rather than passing g around everywhere

	if err := gui.SetColorScheme(); err != nil {
		return err
	}

	if gui.Config.GetUserConfig().GetString("reporting") == "undetermined" {
		gui.waitForIntro.Add(2)
	} else {
		gui.waitForIntro.Add(1)
	}

	if gui.Config.GetUserConfig().GetBool("git.autoFetch") {
		go gui.startBackgroundFetch()
	}
	gui.goEvery(time.Second*10, gui.refreshFiles)
	gui.goEvery(time.Millisecond*50, gui.renderAppStatus)

	g.SetManager(gocui.ManagerFunc(gui.layout), gocui.ManagerFunc(gui.getFocusLayout()))

	if err = gui.keybindings(g); err != nil {
		return err
	}

	err = g.MainLoop()
	return err
}

// RunWithSubprocesses loops, instantiating a new gocui.Gui with each iteration
// if the error returned from a run is a ErrSubProcess, it runs the subprocess
// otherwise it handles the error, possibly by quitting the application
func (gui *Gui) RunWithSubprocesses() error {
	for {
		if err := gui.Run(); err != nil {
			if err == gocui.ErrQuit {
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

func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	if gui.State.Updating {
		return gui.createUpdateQuitConfirmation(g, v)
	}
	if gui.Config.GetUserConfig().GetBool("confirmOnQuit") {
		return gui.createConfirmationPanel(g, v, "", gui.Tr.SLocalize("ConfirmQuit"), func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}, nil)
	}
	return gocui.ErrQuit
}

func (gui *Gui) handleDonate(g *gocui.Gui, v *gocui.View) error {
	if !gui.g.Mouse {
		return nil
	}

	cx, _ := v.Cursor()
	if cx > len(gui.Tr.SLocalize("Donate")) {
		return nil
	}
	return gui.OSCommand.OpenLink("https://donorbox.org/lazygit")
}
