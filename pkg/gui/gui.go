package gui

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"

	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/lbl"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/filtering"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

const StartupPopupVersion = 5

// OverlappingEdges determines if panel edges overlap
var OverlappingEdges = false

func init() {
	runewidth.DefaultCondition.EastAsianWidth = false
}

type ContextManager struct {
	ContextStack []Context
	sync.RWMutex
}

func NewContextManager(initialContext Context) ContextManager {
	return ContextManager{
		ContextStack: []Context{initialContext},
		RWMutex:      sync.RWMutex{},
	}
}

type Repo string

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	g   *gocui.Gui
	Log *logrus.Entry
	Git *commands.Git
	OS  *oscommands.OS

	// this is the state of the GUI for the current repo
	State *guiState

	// this is a mapping of repos to gui states, so that we can restore the original
	// gui state when returning from a subrepo
	RepoStateMap         map[Repo]*guiState
	Config               config.AppConfigurer
	Tr                   *i18n.TranslationSet
	Updater              *updates.Updater
	statusManager        *statusManager
	credentials          credentials
	waitForIntro         sync.WaitGroup
	fileWatcher          *fileWatcher
	viewBufferManagerMap map[string]*tasks.ViewBufferManager
	stopChan             chan struct{}

	// when lazygit is opened outside a git directory we want to open to the most
	// recent repo with the recent repos popup showing
	showRecentRepos bool

	Mutexes guiMutexes

	// findSuggestions will take a string that the user has typed into a prompt
	// and return a slice of suggestions which match that string.
	findSuggestions func(string) []*Suggestion

	// when you enter into a submodule we'll append the superproject's path to this array
	// so that you can return to the superproject
	RepoPathStack []string

	// this tells us whether our views have been initially set up
	ViewsSetup bool

	Views Views

	// if we've suspended the gui (e.g. because we've switched to a subprocess)
	// we typically want to pause some things that are running like background
	// file refreshes
	PauseBackgroundThreads bool

	// Log of the commands that get run, to be displayed to the user.
	CmdLog       []string
	OnRunCommand func(entry oscommands.CmdLogEntry)

	// the extras window contains things like the command log
	ShowExtrasWindow bool

	// used to keep track of which popup is currently shown, so that if we want to
	// open a loading popup for a brief period of time, we can ensure to only remove
	// it if the id hasn't been incremented
	PopupPanelId int
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
type LblPanelState struct {
	*lbl.State
	SecondaryFocused bool // this is for if we show the left or right panel
}

type MergingPanelState struct {
	*mergeconflicts.State

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
	LineByLine     *LblPanelState
	Merging        *MergingPanelState
	CommitFiles    *commitFilesPanelState
	Submodules     *submodulePanelState
	Suggestions    *suggestionsPanelState
}

type Views struct {
	Status        *gocui.View
	Files         *gocui.View
	Branches      *gocui.View
	Commits       *gocui.View
	Stash         *gocui.View
	Main          *gocui.View
	Secondary     *gocui.View
	Options       *gocui.View
	Confirmation  *gocui.View
	Menu          *gocui.View
	Credentials   *gocui.View
	CommitMessage *gocui.View
	CommitFiles   *gocui.View
	Information   *gocui.View
	AppStatus     *gocui.View
	Search        *gocui.View
	SearchPrefix  *gocui.View
	Limit         *gocui.View
	Suggestions   *gocui.View
	Extras        *gocui.View
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

type Modes struct {
	Filtering     filtering.Filtering
	CherryPicking cherrypicking.CherryPicking
	Diffing       diffing.Diffing
	PatchManager  *patch.PatchManager
}

type guiMutexes struct {
	RefreshingFilesMutex  sync.Mutex
	RefreshingStatusMutex sync.Mutex
	FetchMutex            sync.Mutex
	BranchCommitsMutex    sync.Mutex
	LineByLinePanelMutex  sync.Mutex
	SubprocessMutex       sync.Mutex
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
	Suggestions []*Suggestion
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
	SplitMainPanel    bool
	MainContext       ContextKey // used to keep the main and secondary views' contexts in sync
	RetainOriginalDir bool
	IsRefreshingFiles bool
	Searching         searchingState
	ScreenMode        WindowMaximisation
	Ptmx              *os.File
	PrevMainWidth     int
	PrevMainHeight    int
	OldInformation    string
	StartupStage      StartupStage // Allows us to not load everything at once

	Modes Modes

	ContextManager    ContextManager
	Contexts          ContextTree
	ViewContextMap    map[string]Context
	ViewTabContextMap map[string][]tabContext

	// WindowViewNameMap is a mapping of windows to the current view of that window.
	// Some views move between windows for example the commitFiles view and when cycling through
	// side windows we need to know which view to give focus to for a given window
	WindowViewNameMap map[string]string

	// tells us whether we've set up our views for the current repo. We'll need to
	// do this whenever we switch back and forth between repos to get the views
	// back in sync with the repo state
	ViewsSetup bool
}

// reuseState determines if we pull the repo state from our repo state map or
// just re-initialize it. For now we're only re-using state when we're going
// in and out of submodules, for the sake of having the cursor back on the submodule
// when we return.
//
// I tried out always reverting to the repo's original state but found that in fact
// it gets a bit confusing to land back in the status panel when visiting a repo
// you've already switched from. There's no doubt some easy way to make the UX
// optimal for all cases but I'm too lazy to think about what that is right now
func (gui *Gui) resetState(filterPath string, reuseState bool) {
	currentDir, err := os.Getwd()

	if reuseState {
		if err == nil {
			if state := gui.RepoStateMap[Repo(currentDir)]; state != nil {
				gui.State = state
				gui.State.ViewsSetup = false
				return
			}
		} else {
			gui.Log.Error(err)
		}
	}

	showTree := gui.Config.GetUserConfig().Gui.ShowFileTree

	contexts := gui.contextTree()

	screenMode := SCREEN_NORMAL
	initialContext := contexts.Files
	if filterPath != "" {
		screenMode = SCREEN_HALF
		initialContext = contexts.BranchCommits
	}

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
			Merging: &MergingPanelState{
				State:         mergeconflicts.NewState(),
				UserScrolling: false,
			},
		},
		Ptmx: nil,
		Modes: Modes{
			Filtering:     filtering.New(filterPath),
			CherryPicking: cherrypicking.New(),
			Diffing:       diffing.New(),
			PatchManager:  gui.Git.NewPatchManager(),
		},
		ViewContextMap:    contexts.initialViewContextMap(),
		ViewTabContextMap: contexts.initialViewTabContextMap(),
		ScreenMode:        screenMode,
		// TODO: put contexts in the context manager
		ContextManager: NewContextManager(initialContext),
		Contexts:       contexts,
	}

	gui.RepoStateMap[Repo(currentDir)] = gui.State
}

// for now the split view will always be on
// NewGui builds a new gui handler
func NewGui(log *logrus.Entry, gitCommand *commands.Git, oSCommand *oscommands.OS, tr *i18n.TranslationSet, config config.AppConfigurer, updater *updates.Updater, filterPath string, showRecentRepos bool) (*Gui, error) {
	gui := &Gui{
		Log:                  log,
		Git:                  gitCommand,
		OS:                   oSCommand,
		Config:               config,
		Tr:                   tr,
		Updater:              updater,
		statusManager:        &statusManager{},
		viewBufferManagerMap: map[string]*tasks.ViewBufferManager{},
		showRecentRepos:      showRecentRepos,
		RepoPathStack:        []string{},
		RepoStateMap:         map[Repo]*guiState{},
		CmdLog:               []string{},
		ShowExtrasWindow:     config.GetUserConfig().Gui.ShowCommandLog,
	}

	gui.Git.SetCredentialHandlers(gui.PromptUserForCredential, gui.handleCredentialError)

	gui.resetState(filterPath, false)

	gui.watchFilesForChanges()

	onRunCommand := gui.GetOnRunCommand()
	oSCommand.SetOnRunCommand(onRunCommand)
	gui.OnRunCommand = onRunCommand

	return gui, nil
}

// Run setup the gui with keybindings and start the mainloop
func (gui *Gui) Run() error {
	recordEvents := recordingEvents()
	playMode := gocui.NORMAL
	if recordEvents {
		playMode = gocui.RECORDING
	} else if replaying() {
		playMode = gocui.REPLAYING
	}

	g, err := gocui.NewGui(gocui.OutputTrue, OverlappingEdges, playMode, headless())
	if err != nil {
		return err
	}

	gui.g = g // TODO: always use gui.g rather than passing g around everywhere
	defer g.Close()

	if replaying() {
		g.RecordingConfig = gocui.RecordingConfig{
			Speed:  getRecordingSpeed(),
			Leeway: 100,
		}

		g.Recording, err = gui.loadRecording()
		if err != nil {
			return err
		}

		go utils.Safe(func() {
			time.Sleep(time.Second * 40)
			log.Fatal("40 seconds is up, lazygit recording took too long to complete")
		})
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

// RunAndHandleError
func (gui *Gui) RunAndHandleError() error {
	gui.stopChan = make(chan struct{})
	return utils.SafeWithError(func() error {
		if err := gui.Run(); err != nil {
			for _, manager := range gui.viewBufferManagerMap {
				manager.Close()
			}

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

				if err := gui.saveRecording(gui.g.Recording); err != nil {
					return err
				}

				return nil

			default:
				return err
			}
		}

		return nil
	})
}

// returns whether command exited without error or not
func (gui *Gui) runSubprocessWithSuspenseAndRefresh(cmdObj ICmdObj) error {
	_, err := gui.runSubprocessWithSuspense(cmdObj)
	if err != nil {
		return err
	}

	if err := gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC}); err != nil {
		return err
	}

	return nil
}

// returns whether command exited without error or not
func (gui *Gui) runSubprocessWithSuspense(cmdObj ICmdObj) (bool, error) {
	gui.Mutexes.SubprocessMutex.Lock()
	defer gui.Mutexes.SubprocessMutex.Unlock()

	if replaying() {
		// we do not yet support running subprocesses within integration tests. So if
		// we're replaying an integration test and we're inside this method, something
		// has gone wrong, so we should fail

		log.Fatal("opening subprocesses not yet supported in integration tests. Chances are that this test is running too fast and a subprocess is accidentally opened")
	}

	if err := gui.g.Suspend(); err != nil {
		return false, gui.SurfaceError(err)
	}

	gui.PauseBackgroundThreads = true

	cmdErr := gui.runSubprocess(cmdObj)

	if err := gui.g.Resume(); err != nil {
		return false, err
	}

	gui.PauseBackgroundThreads = false

	return cmdErr == nil, gui.SurfaceError(cmdErr)
}

func (gui *Gui) runSubprocess(cmdObj ICmdObj) error {
	cmd := cmdObj.GetCmd()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin

	fmt.Fprintf(os.Stdout, "\n%s\n\n", utils.ColoredString("+ "+strings.Join(cmd.Args, " "), color.FgBlue))

	if err := cmd.Run(); err != nil {
		// not handling the error explicitly because usually we're going to see it
		// in the output anyway
		gui.Log.Error(err)
	}

	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	cmd.Stdin = nil

	fmt.Fprintf(os.Stdout, "\n%s", utils.ColoredString(gui.Tr.PressEnterToReturn, color.FgGreen))
	fmt.Scanln() // wait for enter press

	return nil
}

func (gui *Gui) loadNewRepo() error {
	if err := gui.updateRecentRepoList(); err != nil {
		return err
	}

	if err := gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC}); err != nil {
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
					_ = gui.SurfaceError(err)
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

	return gui.Ask(AskOpts{
		Title:         "",
		Prompt:        gui.Tr.IntroPopupMessage,
		HandleConfirm: onConfirm,
		HandleClose:   onConfirm,
	})
}

func (gui *Gui) goEvery(interval time.Duration, stop chan struct{}, function func() error) {
	go utils.Safe(func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if gui.PauseBackgroundThreads {
					continue
				}
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
	err := gui.fetchInBackground()
	if err != nil && strings.Contains(err.Error(), "exit status 128") && isNew {
		_ = gui.Ask(AskOpts{
			Title:  gui.Tr.NoAutomaticGitFetchTitle,
			Prompt: gui.Tr.NoAutomaticGitFetchBody,
		})
	} else {
		gui.goEvery(time.Second*time.Duration(userConfig.Refresher.FetchInterval), gui.stopChan, func() error {
			err := gui.fetchInBackground()
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
	gui.g.FrameColor = theme.InactiveBorderColor
	gui.g.SelFrameColor = theme.ActiveBorderColor

	return nil
}

func (gui *Gui) GetUserConfig() *config.UserConfig {
	return gui.Config.GetUserConfig()
}

func (gui *Gui) GetGit() commands.IGit {
	return gui.Git
}

func (gui *Gui) GetOS() *oscommands.OS {
	return gui.OS
}

func (gui *Gui) GetTr() *i18n.TranslationSet {
	return gui.Tr
}
