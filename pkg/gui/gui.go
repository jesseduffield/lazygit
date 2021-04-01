package gui

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"

	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/jesseduffield/termbox-go"
	"github.com/mattn/go-runewidth"
	"github.com/sirupsen/logrus"
)

// screen sizing determines how much space your selected window takes up (window
// as in panel, not your terminal's window). Sometimes you want a bit more space
// to see the contents of a panel, and this keeps track of how much maximisation
// you've set
type WindowMaximisation int

const (
	SCREEN_NORMAL WindowMaximisation = iota
	SCREEN_HALF
	SCREEN_FULL
)

const StartupPopupVersion = 3

// OverlappingEdges determines if panel edges overlap
var OverlappingEdges = false

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	g                    *gocui.Gui
	Log                  *logrus.Entry
	GitCommand           *commands.GitCommand
	OSCommand            *oscommands.OSCommand
	SubProcess           *exec.Cmd
	State                *guiState
	Config               config.AppConfigurer
	Tr                   *i18n.TranslationSet
	Errors               SentinelErrors
	Updater              *updates.Updater
	statusManager        *statusManager
	credentials          credentials
	waitForIntro         sync.WaitGroup
	fileWatcher          *fileWatcher
	viewBufferManagerMap map[string]*tasks.ViewBufferManager
	stopChan             chan struct{}

	// when lazygit is opened outside a git directory we want to open to the most
	// recent repo with the recent repos popup showing
	showRecentRepos   bool
	Contexts          ContextTree
	ViewTabContextMap map[string][]tabContext

	// this array either includes the events that we're recording in this session
	// or the events we've recorded in a prior session
	RecordedEvents []RecordedEvent
	StartTime      time.Time

	Mutexes guiStateMutexes

	// findSuggestions will take a string that the user has typed into a prompt
	// and return a slice of suggestions which match that string.
	findSuggestions func(string) []*types.Suggestion
}

type RecordedEvent struct {
	Timestamp int64
	Event     *termbox.Event
}

type listPanelState struct {
	SelectedLineIdx int
}

func (h *listPanelState) SetSelectedLineIdx(value int) {
	h.SelectedLineIdx = value
}

func (h *listPanelState) GetSelectedLineIdx() int {
	return h.SelectedLineIdx
}

// for now the staging panel state, unlike the other panel states, is going to be
// non-mutative, so that we don't accidentally end up
// with mismatches of data. We might change this in the future
type lBlPanelState struct {
	SelectedLineIdx  int
	FirstLineIdx     int
	LastLineIdx      int
	Diff             string
	PatchParser      *patch.PatchParser
	SelectMode       SelectMode
	SecondaryFocused bool // this is for if we show the left or right panel
}

type mergingPanelState struct {
	ConflictIndex  int
	ConflictTop    bool
	Conflicts      []commands.Conflict
	ConflictsMutex sync.Mutex
	EditHistory    *stack.Stack

	// UserScrolling tells us if the user has started scrolling through the file themselves
	// in which case we won't auto-scroll to a conflict.
	UserScrolling bool
}

type filePanelState struct {
	listPanelState
}

// TODO: consider splitting this out into the window and the branches view
type branchPanelState struct {
	listPanelState
}

type remotePanelState struct {
	listPanelState
}

type remoteBranchesState struct {
	listPanelState
}

type tagsPanelState struct {
	listPanelState
}

type commitPanelState struct {
	listPanelState

	LimitCommits bool
}

type reflogCommitPanelState struct {
	listPanelState
}

type subCommitPanelState struct {
	listPanelState

	// e.g. name of branch whose commits we're looking at
	refName string
}

type stashPanelState struct {
	listPanelState
}

type menuPanelState struct {
	listPanelState
	OnPress func() error
}

type commitFilesPanelState struct {
	listPanelState

	// this is the SHA of the commit or the stash index of the stash.
	// Not sure if ref is actually the right word here
	refName   string
	canRebase bool
}

type submodulePanelState struct {
	listPanelState
}

type suggestionsPanelState struct {
	listPanelState
}

type panelStates struct {
	Files          *filePanelState
	Branches       *branchPanelState
	Remotes        *remotePanelState
	RemoteBranches *remoteBranchesState
	Tags           *tagsPanelState
	Commits        *commitPanelState
	ReflogCommits  *reflogCommitPanelState
	SubCommits     *subCommitPanelState
	Stash          *stashPanelState
	Menu           *menuPanelState
	LineByLine     *lBlPanelState
	Merging        *mergingPanelState
	CommitFiles    *commitFilesPanelState
	Submodules     *submodulePanelState
	Suggestions    *suggestionsPanelState
}

type searchingState struct {
	view         *gocui.View
	isSearching  bool
	searchString string
}

// startup stages so we don't need to load everything at once
type StartupStage int

const (
	INITIAL StartupStage = iota
	COMPLETE
)

// if ref is blank we're not diffing anything
type Diffing struct {
	Ref     string
	Reverse bool
}

func (m *Diffing) Active() bool {
	return m.Ref != ""
}

type Filtering struct {
	Path string // the filename that gets passed to git log
}

func (m *Filtering) Active() bool {
	return m.Path != ""
}

type CherryPicking struct {
	CherryPickedCommits []*models.Commit

	// we only allow cherry picking from one context at a time, so you can't copy a commit from the local commits context and then also copy a commit in the reflog context
	ContextKey string
}

func (m *CherryPicking) Active() bool {
	return len(m.CherryPickedCommits) > 0
}

type Modes struct {
	Filtering     Filtering
	CherryPicking CherryPicking
	Diffing       Diffing
}

type guiStateMutexes struct {
	RefreshingFilesMutex  sync.Mutex
	RefreshingStatusMutex sync.Mutex
	FetchMutex            sync.Mutex
	BranchCommitsMutex    sync.Mutex
	LineByLinePanelMutex  sync.Mutex
}

type guiState struct {
	// the file panels (files and commit files) can render as a tree, so we have
	// managers for them which handle rendering a flat list of files in tree form
	FileManager       *filetree.FileManager
	CommitFileManager *filetree.CommitFileManager
	Submodules        []*models.SubmoduleConfig
	Branches          []*models.Branch
	Commits           []*models.Commit
	StashEntries      []*models.StashEntry
	// Suggestions will sometimes appear when typing into a prompt
	Suggestions []*types.Suggestion
	// FilteredReflogCommits are the ones that appear in the reflog panel.
	// when in filtering mode we only include the ones that match the given path
	FilteredReflogCommits []*models.Commit
	// ReflogCommits are the ones used by the branches panel to obtain recency values
	// if we're not in filtering mode, CommitFiles and FilteredReflogCommits will be
	// one and the same
	ReflogCommits     []*models.Commit
	SubCommits        []*models.Commit
	Remotes           []*models.Remote
	RemoteBranches    []*models.RemoteBranch
	Tags              []*models.Tag
	MenuItems         []*menuItem
	Updating          bool
	Panels            *panelStates
	MainContext       string // used to keep the main and secondary views' contexts in sync
	SplitMainPanel    bool
	RetainOriginalDir bool
	IsRefreshingFiles bool
	Searching         searchingState
	ScreenMode        WindowMaximisation
	SideView          *gocui.View
	Ptmx              *os.File
	PrevMainWidth     int
	PrevMainHeight    int
	OldInformation    string
	StartupStage      StartupStage // Allows us to not load everything at once

	Modes Modes

	ContextStack   []Context
	ViewContextMap map[string]Context

	// WindowViewNameMap is a mapping of windows to the current view of that window.
	// Some views move between windows for example the commitFiles view and when cycling through
	// side windows we need to know which view to give focus to for a given window
	WindowViewNameMap map[string]string

	// when you enter into a submodule we'll append the superproject's path to this array
	// so that you can return to the superproject
	RepoPathStack []string
}

func (gui *Gui) resetState() {
	// we carry over the filter path and diff state
	prevFiltering := Filtering{
		Path: "",
	}
	prevDiff := Diffing{}
	prevCherryPicking := CherryPicking{
		CherryPickedCommits: make([]*models.Commit, 0),
		ContextKey:          "",
	}
	prevRepoPathStack := []string{}
	if gui.State != nil {
		prevFiltering = gui.State.Modes.Filtering
		prevDiff = gui.State.Modes.Diffing
		prevCherryPicking = gui.State.Modes.CherryPicking
		prevRepoPathStack = gui.State.RepoPathStack
	}

	modes := Modes{
		Filtering:     prevFiltering,
		CherryPicking: prevCherryPicking,
		Diffing:       prevDiff,
	}

	showTree := gui.Config.GetUserConfig().Gui.ShowFileTree

	gui.State = &guiState{
		FileManager:           filetree.NewFileManager(make([]*models.File, 0), gui.Log, showTree),
		CommitFileManager:     filetree.NewCommitFileManager(make([]*models.CommitFile, 0), gui.Log, showTree),
		Commits:               make([]*models.Commit, 0),
		FilteredReflogCommits: make([]*models.Commit, 0),
		ReflogCommits:         make([]*models.Commit, 0),
		StashEntries:          make([]*models.StashEntry, 0),
		Panels: &panelStates{
			// TODO: work out why some of these are -1 and some are 0. Last time I checked there was a good reason but I'm less certain now
			Files:          &filePanelState{listPanelState{SelectedLineIdx: -1}},
			Submodules:     &submodulePanelState{listPanelState{SelectedLineIdx: -1}},
			Branches:       &branchPanelState{listPanelState{SelectedLineIdx: 0}},
			Remotes:        &remotePanelState{listPanelState{SelectedLineIdx: 0}},
			RemoteBranches: &remoteBranchesState{listPanelState{SelectedLineIdx: -1}},
			Tags:           &tagsPanelState{listPanelState{SelectedLineIdx: -1}},
			Commits:        &commitPanelState{listPanelState: listPanelState{SelectedLineIdx: -1}, LimitCommits: true},
			ReflogCommits:  &reflogCommitPanelState{listPanelState{SelectedLineIdx: 0}},
			SubCommits:     &subCommitPanelState{listPanelState: listPanelState{SelectedLineIdx: 0}, refName: ""},
			CommitFiles:    &commitFilesPanelState{listPanelState: listPanelState{SelectedLineIdx: -1}, refName: ""},
			Stash:          &stashPanelState{listPanelState{SelectedLineIdx: -1}},
			Menu:           &menuPanelState{listPanelState: listPanelState{SelectedLineIdx: 0}, OnPress: nil},
			Suggestions:    &suggestionsPanelState{listPanelState: listPanelState{SelectedLineIdx: 0}},
			Merging: &mergingPanelState{
				ConflictIndex:  0,
				ConflictTop:    true,
				Conflicts:      []commands.Conflict{},
				EditHistory:    stack.New(),
				ConflictsMutex: sync.Mutex{},
			},
		},
		SideView:       nil,
		Ptmx:           nil,
		Modes:          modes,
		ViewContextMap: gui.initialViewContextMap(),
		RepoPathStack:  prevRepoPathStack,
	}
}

// for now the split view will always be on
// NewGui builds a new gui handler
func NewGui(log *logrus.Entry, gitCommand *commands.GitCommand, oSCommand *oscommands.OSCommand, tr *i18n.TranslationSet, config config.AppConfigurer, updater *updates.Updater, filterPath string, showRecentRepos bool) (*Gui, error) {
	gui := &Gui{
		Log:                  log,
		GitCommand:           gitCommand,
		OSCommand:            oSCommand,
		Config:               config,
		Tr:                   tr,
		Updater:              updater,
		statusManager:        &statusManager{},
		viewBufferManagerMap: map[string]*tasks.ViewBufferManager{},
		showRecentRepos:      showRecentRepos,
		RecordedEvents:       []RecordedEvent{},
	}

	gui.resetState()
	gui.State.Modes.Filtering.Path = filterPath
	gui.Contexts = gui.contextTree()
	gui.ViewTabContextMap = gui.viewTabContextMap()

	gui.watchFilesForChanges()

	gui.GenerateSentinelErrors()

	return gui, nil
}

// Run setup the gui with keybindings and start the mainloop
func (gui *Gui) Run() error {
	gui.resetState()

	recordEvents := recordingEvents()

	g, err := gocui.NewGui(gocui.Output256, OverlappingEdges, recordEvents)
	if err != nil {
		return err
	}
	gui.g = g // TODO: always use gui.g rather than passing g around everywhere
	defer g.Close()

	if recordEvents {
		go utils.Safe(gui.recordEvents)
	}

	if gui.State.Modes.Filtering.Active() {
		gui.State.ScreenMode = SCREEN_HALF
	} else {
		gui.State.ScreenMode = SCREEN_NORMAL
	}

	g.OnSearchEscape = gui.onSearchEscape
	if err := gui.Config.ReloadUserConfig(); err != nil {
		return nil
	}
	userConfig := gui.Config.GetUserConfig()
	g.SearchEscapeKey = gui.getKey(userConfig.Keybinding.Universal.Return)
	g.NextSearchMatchKey = gui.getKey(userConfig.Keybinding.Universal.NextMatch)
	g.PrevSearchMatchKey = gui.getKey(userConfig.Keybinding.Universal.PrevMatch)

	g.ASCII = runtime.GOOS == "windows" && runewidth.IsEastAsian()

	if userConfig.Gui.MouseEvents {
		g.Mouse = true
	}

	if err := gui.setColorScheme(); err != nil {
		return err
	}

	if !gui.Config.GetUserConfig().DisableStartupPopups {
		popupTasks := []func(chan struct{}) error{}
		storedPopupVersion := gui.Config.GetAppState().StartupPopupVersion
		if storedPopupVersion < StartupPopupVersion {
			popupTasks = append(popupTasks, gui.showIntroPopupMessage)
		}
		gui.showInitialPopups(popupTasks)
	}

	gui.waitForIntro.Add(1)
	if gui.Config.GetUserConfig().Git.AutoFetch {
		go utils.Safe(gui.startBackgroundFetch)
	}

	gui.goEvery(time.Second*time.Duration(userConfig.Refresher.RefreshInterval), gui.stopChan, gui.refreshFilesAndSubmodules)

	g.SetManager(gocui.ManagerFunc(gui.layout), gocui.ManagerFunc(gui.getFocusLayout()))

	gui.Log.Info("starting main loop")

	err = g.MainLoop()
	return err
}

// RunWithSubprocesses loops, instantiating a new gocui.Gui with each iteration
// if the error returned from a run is a ErrSubProcess, it runs the subprocess
// otherwise it handles the error, possibly by quitting the application
func (gui *Gui) RunWithSubprocesses() error {
	gui.StartTime = time.Now()
	go utils.Safe(gui.replayRecordedEvents)

	for {
		gui.stopChan = make(chan struct{})
		if err := gui.Run(); err != nil {
			for _, manager := range gui.viewBufferManagerMap {
				manager.Close()
			}
			gui.viewBufferManagerMap = map[string]*tasks.ViewBufferManager{}

			if !gui.fileWatcher.Disabled {
				gui.fileWatcher.Watcher.Close()
			}

			close(gui.stopChan)

			switch err {
			case gocui.ErrQuit:
				if !gui.State.RetainOriginalDir {
					if err := gui.recordCurrentDirectory(); err != nil {
						return err
					}
				}

				if err := gui.saveRecordedEvents(); err != nil {
					return err
				}

				return nil
			case gui.Errors.ErrSwitchRepo, gui.Errors.ErrRestart:
				continue
			case gui.Errors.ErrSubProcess:

				if err := gui.runCommand(); err != nil {
					return err
				}
			default:
				return err
			}
		}
	}
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

	fmt.Fprintf(os.Stdout, "\n%s", utils.ColoredString(gui.Tr.PressEnterToReturn, color.FgGreen))
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

	go utils.Safe(func() {
		for _, task := range tasks {
			task := task
			go utils.Safe(func() {
				if err := task(done); err != nil {
					_ = gui.surfaceError(err)
				}
			})

			<-done
			gui.waitForIntro.Done()
		}
	})
}

func (gui *Gui) showIntroPopupMessage(done chan struct{}) error {
	onConfirm := func() error {
		done <- struct{}{}
		gui.Config.GetAppState().StartupPopupVersion = StartupPopupVersion
		return gui.Config.SaveAppState()
	}

	return gui.ask(askOpts{
		title:         "",
		prompt:        gui.Tr.IntroPopupMessage,
		handleConfirm: onConfirm,
		handleClose:   onConfirm,
	})
}

func (gui *Gui) goEvery(interval time.Duration, stop chan struct{}, function func() error) {
	go utils.Safe(func() {
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
	})
}

func (gui *Gui) startBackgroundFetch() {
	gui.waitForIntro.Wait()
	isNew := gui.Config.GetIsNewRepo()
	userConfig := gui.Config.GetUserConfig()
	if !isNew {
		time.After(time.Duration(userConfig.Refresher.FetchInterval) * time.Second)
	}
	err := gui.fetch(false)
	if err != nil && strings.Contains(err.Error(), "exit status 128") && isNew {
		_ = gui.ask(askOpts{
			title:  gui.Tr.NoAutomaticGitFetchTitle,
			prompt: gui.Tr.NoAutomaticGitFetchBody,
		})
	} else {
		gui.goEvery(time.Second*time.Duration(userConfig.Refresher.FetchInterval), gui.stopChan, func() error {
			err := gui.fetch(false)
			return err
		})
	}
}

// setColorScheme sets the color scheme for the app based on the user config
func (gui *Gui) setColorScheme() error {
	userConfig := gui.Config.GetUserConfig()
	theme.UpdateTheme(userConfig.Gui.Theme)

	gui.g.FgColor = theme.InactiveBorderColor
	gui.g.SelFgColor = theme.ActiveBorderColor

	return nil
}
