package gui

import (
	"fmt"
	"io/ioutil"
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
	ErrRestart    error
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
		ErrRestart:    errors.New("restarting"),
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
	State                *guiState
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
	SelectedLine int
	LimitCommits bool
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
}

type searchingState struct {
	view         *gocui.View
	isSearching  bool
	searchString string
}

// startup stages so we don't need to load everything at once
const (
	INITIAL = iota
	COMPLETE
)

// if ref is blank we're not diffing anything
type DiffState struct {
	Ref     string
	Reverse bool
}

type guiState struct {
	Files        []*commands.File
	Branches     []*commands.Branch
	Commits      []*commands.Commit
	StashEntries []*commands.StashEntry
	CommitFiles  []*commands.CommitFile
	// FilteredReflogCommits are the ones that appear in the reflog panel.
	// when in filtering mode we only include the ones that match the given path
	FilteredReflogCommits []*commands.Commit
	// ReflogCommits are the ones used by the branches panel to obtain recency values
	// if we're not in filtering mode, CommitFiles and FilteredReflogCommits will be
	// one and the same
	ReflogCommits         []*commands.Commit
	Remotes               []*commands.Remote
	RemoteBranches        []*commands.RemoteBranch
	Tags                  []*commands.Tag
	MenuItemCount         int // can't store the actual list because it's of interface{} type
	PreviousView          string
	Updating              bool
	Panels                *panelStates
	MainContext           string // used to keep the main and secondary views' contexts in sync
	CherryPickedCommits   []*commands.Commit
	SplitMainPanel        bool
	RetainOriginalDir     bool
	IsRefreshingFiles     bool
	RefreshingFilesMutex  sync.Mutex
	RefreshingStatusMutex sync.Mutex
	Searching             searchingState
	ScreenMode            int
	SideView              *gocui.View
	Ptmx                  *os.File
	PrevMainWidth         int
	PrevMainHeight        int
	OldInformation        string
	StartupStage          int    // one of INITIAL and COMPLETE. Allows us to not load everything at once
	FilterPath            string // the filename that gets passed to git log
	Diff                  DiffState
}

func (gui *Gui) resetState() {
	// we carry over the filter path and diff state
	prevFilterPath := ""
	prevDiff := DiffState{}
	if gui.State != nil {
		prevFilterPath = gui.State.FilterPath
		prevDiff = gui.State.Diff
	}

	gui.State = &guiState{
		Files:                 make([]*commands.File, 0),
		PreviousView:          "files",
		Commits:               make([]*commands.Commit, 0),
		FilteredReflogCommits: make([]*commands.Commit, 0),
		ReflogCommits:         make([]*commands.Commit, 0),
		CherryPickedCommits:   make([]*commands.Commit, 0),
		StashEntries:          make([]*commands.StashEntry, 0),
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
		},
		SideView:   nil,
		Ptmx:       nil,
		FilterPath: prevFilterPath,
		Diff:       prevDiff,
	}
}

// for now the split view will always be on
// NewGui builds a new gui handler
func NewGui(log *logrus.Entry, gitCommand *commands.GitCommand, oSCommand *commands.OSCommand, tr *i18n.Localizer, config config.AppConfigurer, updater *updates.Updater, filterPath string) (*Gui, error) {
	gui := &Gui{
		Log:                  log,
		GitCommand:           gitCommand,
		OSCommand:            oSCommand,
		Config:               config,
		Tr:                   tr,
		Updater:              updater,
		statusManager:        &statusManager{},
		viewBufferManagerMap: map[string]*tasks.ViewBufferManager{},
	}

	gui.resetState()
	gui.State.FilterPath = filterPath

	gui.watchFilesForChanges()

	gui.GenerateSentinelErrors()

	return gui, nil
}

// Run setup the gui with keybindings and start the mainloop
func (gui *Gui) Run() error {
	gui.resetState()

	g, err := gocui.NewGui(gocui.Output256, OverlappingEdges)
	if err != nil {
		return err
	}
	defer g.Close()

	if gui.inFilterMode() {
		gui.State.ScreenMode = SCREEN_HALF
	} else {
		gui.State.ScreenMode = SCREEN_NORMAL
	}

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
			} else if err == gui.Errors.ErrRestart {
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

func (gui *Gui) loadNewRepo() error {
	gui.Updater.CheckForNewUpdate(gui.onBackgroundUpdateCheckFinish, false)
	if err := gui.updateRecentRepoList(); err != nil {
		return err
	}
	gui.waitForIntro.Done()

	if err := gui.refreshSidePanels(refreshOptions{mode: ASYNC}); err != nil {
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
					_ = gui.surfaceError(err)
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

// setColorScheme sets the color scheme for the app based on the user config
func (gui *Gui) setColorScheme() error {
	userConfig := gui.Config.GetUserConfig()
	theme.UpdateTheme(userConfig)

	gui.g.FgColor = theme.InactiveBorderColor
	gui.g.SelFgColor = theme.ActiveBorderColor

	return nil
}
