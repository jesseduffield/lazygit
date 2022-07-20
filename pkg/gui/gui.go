package gui

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/lbl"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/filtering"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/authors"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/graph"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/services/custom_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"gopkg.in/ozeidan/fuzzy-patricia.v3/patricia"
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

type ContextManager struct {
	ContextStack []types.Context
	sync.RWMutex
}

func NewContextManager(initialContext types.Context) ContextManager {
	return ContextManager{
		ContextStack: []types.Context{initialContext},
		RWMutex:      sync.RWMutex{},
	}
}

type Repo string

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	*common.Common
	g   *gocui.Gui
	git *commands.GitCommand
	os  *oscommands.OSCommand

	// this is the state of the GUI for the current repo
	State *GuiRepoState

	CustomCommandsClient *custom_commands.Client

	// this is a mapping of repos to gui states, so that we can restore the original
	// gui state when returning from a subrepo
	RepoStateMap         map[Repo]*GuiRepoState
	Config               config.AppConfigurer
	Updater              *updates.Updater
	statusManager        *statusManager
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
	findSuggestions func(string) []*types.Suggestion

	// when you enter into a submodule we'll append the superproject's path to this array
	// so that you can return to the superproject
	RepoPathStack *utils.StringStack

	// this tells us whether our views have been initially set up
	ViewsSetup bool

	Views Views

	// if we've suspended the gui (e.g. because we've switched to a subprocess)
	// we typically want to pause some things that are running like background
	// file refreshes
	PauseBackgroundThreads bool

	// Log of the commands that get run, to be displayed to the user.
	CmdLog []string

	// the extras window contains things like the command log
	ShowExtrasWindow bool

	suggestionsAsyncHandler *tasks.AsyncHandler

	PopupHandler types.IPopupHandler

	IsNewRepo bool

	// flag as to whether or not the diff view should ignore whitespace
	IgnoreWhitespaceInDiffView bool

	// we use this to decide whether we'll return to the original directory that
	// lazygit was opened in, or if we'll retain the one we're currently in.
	RetainOriginalDir bool

	PrevLayout PrevLayout

	// this is the initial dir we are in upon opening lazygit. We hold onto this
	// in case we want to restore it before quitting for users who have set up
	// the feature for changing directory upon quit.
	// The reason we don't just wait until quit time to handle changing directories
	// is because some users want to keep track of the current lazygit directory in an outside
	// process
	InitialDir string

	c       *types.HelperCommon
	helpers *helpers.Helpers
}

// we keep track of some stuff from one render to the next to see if certain
// things have changed
type PrevLayout struct {
	Information string
	MainWidth   int
	MainHeight  int
}

type GuiRepoState struct {
	Model *types.Model
	Modes *types.Modes

	// Suggestions will sometimes appear when typing into a prompt
	Suggestions []*types.Suggestion

	Updating       bool
	Panels         *panelStates
	SplitMainPanel bool
	LimitCommits   bool

	IsRefreshingFiles bool
	Searching         searchingState
	Ptmx              *os.File
	StartupStage      StartupStage // Allows us to not load everything at once

	MainContext       types.ContextKey // used to keep the main and secondary views' contexts in sync
	ContextManager    ContextManager
	Contexts          *context.ContextTree
	ViewContextMap    *context.ViewContextMap
	ViewTabContextMap map[string][]context.TabContext

	// WindowViewNameMap is a mapping of windows to the current view of that window.
	// Some views move between windows for example the commitFiles view and when cycling through
	// side windows we need to know which view to give focus to for a given window
	WindowViewNameMap map[string]string

	// tells us whether we've set up our views for the current repo. We'll need to
	// do this whenever we switch back and forth between repos to get the views
	// back in sync with the repo state
	ViewsSetup bool

	// we store a commit message in this field if we've escaped the commit message
	// panel without committing or if our commit failed
	savedCommitMessage string

	ScreenMode WindowMaximisation

	CurrentPopupOpts *types.CreatePopupPanelOpts
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

	// UserVerticalScrolling tells us if the user has started scrolling through the file themselves
	// in which case we won't auto-scroll to a conflict.
	UserVerticalScrolling bool
}

// as we move things to the new context approach we're going to eventually
// remove this struct altogether and store this state on the contexts.
type panelStates struct {
	LineByLine *LblPanelState
	Merging    *MergingPanelState
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

// if you add a new mutex here be sure to instantiate it. We're using pointers to
// mutexes so that we can pass the mutexes to controllers.
type guiMutexes struct {
	RefreshingFilesMutex  *sync.Mutex
	RefreshingStatusMutex *sync.Mutex
	SyncMutex             *sync.Mutex
	LocalCommitsMutex     *sync.Mutex
	LineByLinePanelMutex  *sync.Mutex
	SubprocessMutex       *sync.Mutex
	PopupMutex            *sync.Mutex
}

func (gui *Gui) onNewRepo(startArgs types.StartArgs, reuseState bool) error {
	var err error
	gui.git, err = commands.NewGitCommand(
		gui.Common,
		gui.os,
		git_config.NewStdCachedGitConfig(gui.Log),
		gui.Mutexes.SyncMutex,
	)
	if err != nil {
		return err
	}

	gui.resetState(startArgs, reuseState)

	gui.resetControllers()

	if err := gui.resetKeybindings(); err != nil {
		return err
	}

	return nil
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
func (gui *Gui) resetState(startArgs types.StartArgs, reuseState bool) {
	currentDir, err := os.Getwd()

	if reuseState {
		if err == nil {
			if state := gui.RepoStateMap[Repo(currentDir)]; state != nil {
				gui.State = state
				gui.State.ViewsSetup = false

				// setting this to nil so we don't get stuck based on a popup that was
				// previously opened
				gui.Mutexes.PopupMutex.Lock()
				gui.State.CurrentPopupOpts = nil
				gui.Mutexes.PopupMutex.Unlock()

				gui.syncViewContexts()
				return
			}
		} else {
			gui.c.Log.Error(err)
		}
	}

	contextTree := gui.contextTree()

	initialContext := initialContext(contextTree, startArgs)
	initialScreenMode := initialScreenMode(startArgs)

	viewContextMap := context.NewViewContextMap()
	for viewName, context := range initialViewContextMapping(contextTree) {
		viewContextMap.Set(viewName, context)
	}

	gui.State = &GuiRepoState{
		Model: &types.Model{
			CommitFiles:           nil,
			Files:                 make([]*models.File, 0),
			Commits:               make([]*models.Commit, 0),
			StashEntries:          make([]*models.StashEntry, 0),
			FilteredReflogCommits: make([]*models.Commit, 0),
			ReflogCommits:         make([]*models.Commit, 0),
			BisectInfo:            git_commands.NewNullBisectInfo(),
			FilesTrie:             patricia.NewTrie(),
		},

		Panels: &panelStates{
			Merging: &MergingPanelState{
				State:                 mergeconflicts.NewState(),
				UserVerticalScrolling: false,
			},
		},
		Ptmx: nil,
		Modes: &types.Modes{
			Filtering:     filtering.New(startArgs.FilterPath),
			CherryPicking: cherrypicking.New(),
			Diffing:       diffing.New(),
		},
		ViewContextMap:    viewContextMap,
		ViewTabContextMap: gui.initialViewTabContextMap(contextTree),
		ScreenMode:        initialScreenMode,
		// TODO: put contexts in the context manager
		ContextManager: NewContextManager(initialContext),
		Contexts:       contextTree,
	}

	gui.syncViewContexts()

	gui.RepoStateMap[Repo(currentDir)] = gui.State
}

func initialScreenMode(startArgs types.StartArgs) WindowMaximisation {
	if startArgs.FilterPath != "" || startArgs.GitArg != types.GitArgNone {
		return SCREEN_HALF
	} else {
		return SCREEN_NORMAL
	}
}

func initialContext(contextTree *context.ContextTree, startArgs types.StartArgs) types.IListContext {
	var initialContext types.IListContext = contextTree.Files

	if startArgs.FilterPath != "" {
		initialContext = contextTree.LocalCommits
	} else if startArgs.GitArg != types.GitArgNone {
		switch startArgs.GitArg {
		case types.GitArgStatus:
			initialContext = contextTree.Files
		case types.GitArgBranch:
			initialContext = contextTree.Branches
		case types.GitArgLog:
			initialContext = contextTree.LocalCommits
		case types.GitArgStash:
			initialContext = contextTree.Stash
		default:
			panic("unhandled git arg")
		}
	}

	return initialContext
}

func (gui *Gui) syncViewContexts() {
	for viewName, context := range gui.State.ViewContextMap.Entries() {
		view, err := gui.g.View(viewName)
		if err != nil {
			panic(err)
		}
		view.Context = string(context.GetKey())
	}
}

// for now the split view will always be on
// NewGui builds a new gui handler
func NewGui(
	cmn *common.Common,
	config config.AppConfigurer,
	gitConfig git_config.IGitConfig,
	updater *updates.Updater,
	showRecentRepos bool,
	initialDir string,
) (*Gui, error) {
	gui := &Gui{
		Common:                  cmn,
		Config:                  config,
		Updater:                 updater,
		statusManager:           &statusManager{},
		viewBufferManagerMap:    map[string]*tasks.ViewBufferManager{},
		showRecentRepos:         showRecentRepos,
		RepoPathStack:           &utils.StringStack{},
		RepoStateMap:            map[Repo]*GuiRepoState{},
		CmdLog:                  []string{},
		suggestionsAsyncHandler: tasks.NewAsyncHandler(),

		// originally we could only hide the command log permanently via the config
		// but now we do it via state. So we need to still support the config for the
		// sake of backwards compatibility. We're making use of short circuiting here
		ShowExtrasWindow: cmn.UserConfig.Gui.ShowCommandLog && !config.GetAppState().HideCommandLog,
		Mutexes: guiMutexes{
			RefreshingFilesMutex:  &sync.Mutex{},
			RefreshingStatusMutex: &sync.Mutex{},
			SyncMutex:             &sync.Mutex{},
			LocalCommitsMutex:     &sync.Mutex{},
			LineByLinePanelMutex:  &sync.Mutex{},
			SubprocessMutex:       &sync.Mutex{},
			PopupMutex:            &sync.Mutex{},
		},
		InitialDir: initialDir,
	}

	gui.watchFilesForChanges()

	gui.PopupHandler = popup.NewPopupHandler(
		cmn,
		gui.createPopupPanel,
		func() error { return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC}) },
		func() error { return gui.closeConfirmationPrompt(false) },
		gui.createMenu,
		gui.withWaitingStatus,
		gui.toast,
		func() string { return gui.Views.Confirmation.TextArea.GetContent() },
	)

	guiCommon := &guiCommon{gui: gui, IPopupHandler: gui.PopupHandler}
	helperCommon := &types.HelperCommon{IGuiCommon: guiCommon, Common: cmn}

	credentialsHelper := helpers.NewCredentialsHelper(helperCommon)

	guiIO := oscommands.NewGuiIO(
		cmn.Log,
		gui.LogCommand,
		gui.getCmdWriter,
		credentialsHelper.PromptUserForCredential,
	)

	osCommand := oscommands.NewOSCommand(cmn, config, oscommands.GetPlatform(), guiIO)

	gui.os = osCommand

	// storing this stuff on the gui for now to ease refactoring
	// TODO: reset these controllers upon changing repos due to state changing
	gui.c = helperCommon

	authors.SetCustomAuthors(gui.UserConfig.Gui.AuthorColors)
	icons.SetIconEnabled(gui.UserConfig.Gui.ShowIcons)
	presentation.SetCustomBranches(gui.UserConfig.Gui.BranchColors)

	return gui, nil
}

var RuneReplacements = map[rune]string{
	// for the commit graph
	graph.MergeSymbol:  "M",
	graph.CommitSymbol: "o",
}

func (gui *Gui) initGocui() (*gocui.Gui, error) {
	recordEvents := recordingEvents()
	playMode := gocui.NORMAL
	if recordEvents {
		playMode = gocui.RECORDING
	} else if replaying() {
		playMode = gocui.REPLAYING
	}

	g, err := gocui.NewGui(gocui.OutputTrue, OverlappingEdges, playMode, headless(), RuneReplacements)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (gui *Gui) initialViewTabContextMap(contextTree *context.ContextTree) map[string][]context.TabContext {
	return map[string][]context.TabContext{
		"branches": {
			{
				Tab:     gui.c.Tr.LocalBranchesTitle,
				Context: contextTree.Branches,
			},
			{
				Tab:     gui.c.Tr.RemotesTitle,
				Context: contextTree.Remotes,
			},
			{
				Tab:     gui.c.Tr.TagsTitle,
				Context: contextTree.Tags,
			},
		},
		"commits": {
			{
				Tab:     gui.c.Tr.CommitsTitle,
				Context: contextTree.LocalCommits,
			},
			{
				Tab:     gui.c.Tr.ReflogCommitsTitle,
				Context: contextTree.ReflogCommits,
			},
		},
		"files": {
			{
				Tab:     gui.c.Tr.FilesTitle,
				Context: contextTree.Files,
			},
			{
				Tab:     gui.c.Tr.SubmodulesTitle,
				Context: contextTree.Submodules,
			},
		},
	}
}

// Run: setup the gui with keybindings and start the mainloop
func (gui *Gui) Run(startArgs types.StartArgs) error {
	g, err := gui.initGocui()
	if err != nil {
		return err
	}

	gui.g = g
	defer gui.g.Close()

	if replaying() {
		gui.g.RecordingConfig = gocui.RecordingConfig{
			Speed:  getRecordingSpeed(),
			Leeway: 100,
		}

		var err error
		gui.g.Recording, err = gui.loadRecording()
		if err != nil {
			return err
		}

		go utils.Safe(func() {
			time.Sleep(time.Second * 40)
			log.Fatal("40 seconds is up, lazygit recording took too long to complete")
		})
	}

	gui.g.OnSearchEscape = gui.onSearchEscape
	if err := gui.Config.ReloadUserConfig(); err != nil {
		return nil
	}
	userConfig := gui.UserConfig
	gui.g.SearchEscapeKey = gui.getKey(userConfig.Keybinding.Universal.Return)
	gui.g.NextSearchMatchKey = gui.getKey(userConfig.Keybinding.Universal.NextMatch)
	gui.g.PrevSearchMatchKey = gui.getKey(userConfig.Keybinding.Universal.PrevMatch)

	gui.g.ShowListFooter = userConfig.Gui.ShowListFooter

	if userConfig.Gui.MouseEvents {
		gui.g.Mouse = true
	}

	if err := gui.setColorScheme(); err != nil {
		return err
	}

	gui.g.SetManager(gocui.ManagerFunc(gui.layout), gocui.ManagerFunc(gui.getFocusLayout()))

	if err := gui.createAllViews(); err != nil {
		return err
	}

	// onNewRepo must be called after g.SetManager because SetManager deletes keybindings
	if err := gui.onNewRepo(startArgs, false); err != nil {
		return err
	}

	gui.waitForIntro.Add(1)

	if userConfig.Git.AutoFetch {
		fetchInterval := userConfig.Refresher.FetchInterval
		if fetchInterval > 0 {
			go utils.Safe(gui.startBackgroundFetch)
		} else {
			gui.c.Log.Errorf(
				"Value of config option 'refresher.fetchInterval' (%d) is invalid, disabling auto-fetch",
				fetchInterval)
		}
	}

	if userConfig.Git.AutoRefresh {
		refreshInterval := userConfig.Refresher.RefreshInterval
		if refreshInterval > 0 {
			gui.goEvery(time.Second*time.Duration(refreshInterval), gui.stopChan, gui.refreshFilesAndSubmodules)
		} else {
			gui.c.Log.Errorf(
				"Value of config option 'refresher.refreshInterval' (%d) is invalid, disabling auto-refresh",
				refreshInterval)
		}
	}

	gui.c.Log.Info("starting main loop")

	return gui.g.MainLoop()
}

func (gui *Gui) RunAndHandleError(startArgs types.StartArgs) error {
	gui.stopChan = make(chan struct{})
	return utils.SafeWithError(func() error {
		if err := gui.Run(startArgs); err != nil {
			for _, manager := range gui.viewBufferManagerMap {
				manager.Close()
			}

			if !gui.fileWatcher.Disabled {
				gui.fileWatcher.Watcher.Close()
			}

			close(gui.stopChan)

			switch err {
			case gocui.ErrQuit:
				if gui.RetainOriginalDir {
					if err := gui.recordDirectory(gui.InitialDir); err != nil {
						return err
					}
				} else {
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
func (gui *Gui) runSubprocessWithSuspenseAndRefresh(subprocess oscommands.ICmdObj) error {
	_, err := gui.runSubprocessWithSuspense(subprocess)
	if err != nil {
		return err
	}

	if err := gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC}); err != nil {
		return err
	}

	return nil
}

// returns whether command exited without error or not
func (gui *Gui) runSubprocessWithSuspense(subprocess oscommands.ICmdObj) (bool, error) {
	gui.Mutexes.SubprocessMutex.Lock()
	defer gui.Mutexes.SubprocessMutex.Unlock()

	if replaying() {
		// we do not yet support running subprocesses within integration tests. So if
		// we're replaying an integration test and we're inside this method, something
		// has gone wrong, so we should fail

		log.Fatal("opening subprocesses not yet supported in integration tests. Chances are that this test is running too fast and a subprocess is accidentally opened")
	}

	if err := gui.g.Suspend(); err != nil {
		return false, gui.c.Error(err)
	}

	gui.PauseBackgroundThreads = true

	cmdErr := gui.runSubprocess(subprocess)

	if err := gui.g.Resume(); err != nil {
		return false, err
	}

	gui.PauseBackgroundThreads = false

	if cmdErr != nil {
		return false, gui.c.Error(cmdErr)
	}

	return true, nil
}

func (gui *Gui) runSubprocess(cmdObj oscommands.ICmdObj) error { //nolint:unparam
	gui.LogCommand(cmdObj.ToString(), true)

	subprocess := cmdObj.GetCmd()
	subprocess.Stdout = os.Stdout
	subprocess.Stderr = os.Stdout
	subprocess.Stdin = os.Stdin

	fmt.Fprintf(os.Stdout, "\n%s\n\n", style.FgBlue.Sprint("+ "+strings.Join(subprocess.Args, " ")))

	err := subprocess.Run()

	subprocess.Stdout = ioutil.Discard
	subprocess.Stderr = ioutil.Discard
	subprocess.Stdin = nil

	if gui.Config.GetUserConfig().PromptToReturnFromSubprocess {
		fmt.Fprintf(os.Stdout, "\n%s", style.FgGreen.Sprint(gui.Tr.PressEnterToReturn))
		fmt.Scanln() // wait for enter press
	}

	return err
}

func (gui *Gui) loadNewRepo() error {
	if err := gui.updateRecentRepoList(); err != nil {
		return err
	}

	if err := gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC}); err != nil {
		return err
	}

	if err := gui.os.UpdateWindowTitle(); err != nil {
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
					_ = gui.c.Error(err)
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
		gui.c.GetAppState().StartupPopupVersion = StartupPopupVersion
		return gui.c.SaveAppState()
	}

	return gui.c.Confirm(types.ConfirmOpts{
		Title:         "",
		Prompt:        gui.c.Tr.IntroPopupMessage,
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
	isNew := gui.IsNewRepo
	userConfig := gui.UserConfig
	if !isNew {
		time.After(time.Duration(userConfig.Refresher.FetchInterval) * time.Second)
	}
	err := gui.backgroundFetch()
	if err != nil && strings.Contains(err.Error(), "exit status 128") && isNew {
		_ = gui.c.Alert(gui.c.Tr.NoAutomaticGitFetchTitle, gui.c.Tr.NoAutomaticGitFetchBody)
	} else {
		gui.goEvery(time.Second*time.Duration(userConfig.Refresher.FetchInterval), gui.stopChan, func() error {
			err := gui.backgroundFetch()
			gui.render()
			return err
		})
	}
}

// setColorScheme sets the color scheme for the app based on the user config
func (gui *Gui) setColorScheme() error {
	userConfig := gui.UserConfig
	theme.UpdateTheme(userConfig.Gui.Theme)

	gui.g.FgColor = theme.InactiveBorderColor
	gui.g.SelFgColor = theme.ActiveBorderColor
	gui.g.FrameColor = theme.InactiveBorderColor
	gui.g.SelFrameColor = theme.ActiveBorderColor

	return nil
}

func (gui *Gui) OnUIThread(f func() error) {
	gui.g.Update(func(*gocui.Gui) error {
		return f()
	})
}
