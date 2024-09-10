package gui

import (
	goContext "context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazycore/pkg/boxlayout"
	appTypes "github.com/jesseduffield/lazygit/pkg/app/types"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/filtering"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/marked_base_commit"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/authors"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/graph"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/services/custom_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/status"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/sasha-s/go-deadlock"
	"gopkg.in/ozeidan/fuzzy-patricia.v3/patricia"
)

const StartupPopupVersion = 5

// OverlappingEdges determines if panel edges overlap
var OverlappingEdges = false

type Repo string

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	*common.Common
	g          *gocui.Gui
	gitVersion *git_commands.GitVersion
	git        *commands.GitCommand
	os         *oscommands.OSCommand

	// this is the state of the GUI for the current repo
	State *GuiRepoState

	CustomCommandsClient *custom_commands.Client

	// this is a mapping of repos to gui states, so that we can restore the original
	// gui state when returning from a subrepo.
	// In repos with multiple worktrees, we store a separate repo state per worktree.
	RepoStateMap         map[Repo]*GuiRepoState
	Config               config.AppConfigurer
	Updater              *updates.Updater
	statusManager        *status.StatusManager
	waitForIntro         sync.WaitGroup
	viewBufferManagerMap map[string]*tasks.ViewBufferManager
	// holds a mapping of view names to ptmx's. This is for rendering command outputs
	// from within a pty. The point of keeping track of them is so that if we re-size
	// the window, we can tell the pty it needs to resize accordingly.
	viewPtmxMap map[string]*os.File
	stopChan    chan struct{}

	// when lazygit is opened outside a git directory we want to open to the most
	// recent repo with the recent repos popup showing
	showRecentRepos bool

	Mutexes types.Mutexes

	// when you enter into a submodule we'll append the superproject's path to this array
	// so that you can return to the superproject
	RepoPathStack *utils.StringStack

	// this tells us whether our views have been initially set up
	ViewsSetup bool

	Views types.Views

	// Log of the commands/actions logged in the Command Log panel.
	GuiLog []string

	// the extras window contains things like the command log
	ShowExtrasWindow bool

	PopupHandler types.IPopupHandler

	IsNewRepo bool

	IsRefreshingFiles bool

	// we use this to decide whether we'll return to the original directory that
	// lazygit was opened in, or if we'll retain the one we're currently in.
	RetainOriginalDir bool

	// stores long-running operations associated with items (e.g. when a branch
	// is being pushed). At the moment the rule is to use an item operation when
	// we need to talk to the remote.
	itemOperations      map[string]types.ItemOperation
	itemOperationsMutex *deadlock.Mutex

	PrevLayout PrevLayout

	// this is the initial dir we are in upon opening lazygit. We hold onto this
	// in case we want to restore it before quitting for users who have set up
	// the feature for changing directory upon quit.
	// The reason we don't just wait until quit time to handle changing directories
	// is because some users want to keep track of the current lazygit directory in an outside
	// process
	InitialDir string

	BackgroundRoutineMgr *BackgroundRoutineMgr
	// for accessing the gui's state from outside this package
	stateAccessor *StateAccessor

	Updating bool

	c       *helpers.HelperCommon
	helpers *helpers.Helpers

	previousLanguageConfig string

	integrationTest integrationTypes.IntegrationTest

	afterLayoutFuncs chan func() error
}

type StateAccessor struct {
	gui *Gui
}

var _ types.IStateAccessor = new(StateAccessor)

func (self *StateAccessor) GetRepoPathStack() *utils.StringStack {
	return self.gui.RepoPathStack
}

func (self *StateAccessor) GetUpdating() bool {
	return self.gui.Updating
}

func (self *StateAccessor) SetUpdating(value bool) {
	self.gui.Updating = value
}

func (self *StateAccessor) GetRepoState() types.IRepoStateAccessor {
	return self.gui.State
}

func (self *StateAccessor) GetIsRefreshingFiles() bool {
	return self.gui.IsRefreshingFiles
}

func (self *StateAccessor) SetIsRefreshingFiles(value bool) {
	self.gui.IsRefreshingFiles = value
}

func (self *StateAccessor) GetShowExtrasWindow() bool {
	return self.gui.ShowExtrasWindow
}

func (self *StateAccessor) SetShowExtrasWindow(value bool) {
	self.gui.ShowExtrasWindow = value
}

func (self *StateAccessor) GetRetainOriginalDir() bool {
	return self.gui.RetainOriginalDir
}

func (self *StateAccessor) SetRetainOriginalDir(value bool) {
	self.gui.RetainOriginalDir = value
}

func (self *StateAccessor) GetItemOperation(item types.HasUrn) types.ItemOperation {
	self.gui.itemOperationsMutex.Lock()
	defer self.gui.itemOperationsMutex.Unlock()

	return self.gui.itemOperations[item.URN()]
}

func (self *StateAccessor) SetItemOperation(item types.HasUrn, operation types.ItemOperation) {
	self.gui.itemOperationsMutex.Lock()
	defer self.gui.itemOperationsMutex.Unlock()

	self.gui.itemOperations[item.URN()] = operation
}

func (self *StateAccessor) ClearItemOperation(item types.HasUrn) {
	self.gui.itemOperationsMutex.Lock()
	defer self.gui.itemOperationsMutex.Unlock()

	delete(self.gui.itemOperations, item.URN())
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

	SplitMainPanel bool
	LimitCommits   bool

	SearchState  *types.SearchState
	StartupStage types.StartupStage // Allows us to not load everything at once

	ContextMgr *ContextMgr
	Contexts   *context.ContextTree

	// WindowViewNameMap is a mapping of windows to the current view of that window.
	// Some views move between windows for example the commitFiles view and when cycling through
	// side windows we need to know which view to give focus to for a given window
	WindowViewNameMap *utils.ThreadSafeMap[string, string]

	// tells us whether we've set up our views for the current repo. We'll need to
	// do this whenever we switch back and forth between repos to get the views
	// back in sync with the repo state
	ViewsSetup bool

	ScreenMode types.WindowMaximisation

	CurrentPopupOpts *types.CreatePopupPanelOpts
}

var _ types.IRepoStateAccessor = new(GuiRepoState)

func (self *GuiRepoState) GetViewsSetup() bool {
	return self.ViewsSetup
}

func (self *GuiRepoState) GetWindowViewNameMap() *utils.ThreadSafeMap[string, string] {
	return self.WindowViewNameMap
}

func (self *GuiRepoState) GetStartupStage() types.StartupStage {
	return self.StartupStage
}

func (self *GuiRepoState) SetStartupStage(value types.StartupStage) {
	self.StartupStage = value
}

func (self *GuiRepoState) GetCurrentPopupOpts() *types.CreatePopupPanelOpts {
	return self.CurrentPopupOpts
}

func (self *GuiRepoState) SetCurrentPopupOpts(value *types.CreatePopupPanelOpts) {
	self.CurrentPopupOpts = value
}

func (self *GuiRepoState) GetScreenMode() types.WindowMaximisation {
	return self.ScreenMode
}

func (self *GuiRepoState) SetScreenMode(value types.WindowMaximisation) {
	self.ScreenMode = value
}

func (self *GuiRepoState) InSearchPrompt() bool {
	return self.SearchState.SearchType() != types.SearchTypeNone
}

func (self *GuiRepoState) GetSearchState() *types.SearchState {
	return self.SearchState
}

func (self *GuiRepoState) SetSplitMainPanel(value bool) {
	self.SplitMainPanel = value
}

func (self *GuiRepoState) GetSplitMainPanel() bool {
	return self.SplitMainPanel
}

func (gui *Gui) onNewRepo(startArgs appTypes.StartArgs, contextKey types.ContextKey) error {
	var err error
	gui.git, err = commands.NewGitCommand(
		gui.Common,
		gui.gitVersion,
		gui.os,
		git_config.NewStdCachedGitConfig(gui.Log),
	)
	if err != nil {
		return err
	}

	err = gui.Config.ReloadUserConfigForRepo(gui.getPerRepoConfigFiles())
	if err != nil {
		return err
	}

	err = gui.onUserConfigLoaded()
	if err != nil {
		return err
	}

	contextToPush := gui.resetState(startArgs)

	gui.resetHelpersAndControllers()

	if err := gui.resetKeybindings(); err != nil {
		return err
	}

	gui.g.SetFocusHandler(func(Focused bool) error {
		if Focused {
			oldConfig := gui.Config.GetUserConfig()
			reloadErr, didChange := gui.Config.ReloadChangedUserConfigFiles()
			if didChange && reloadErr == nil {
				gui.c.Log.Info("User config changed - reloading")
				reloadErr = gui.onUserConfigLoaded()
				if err := gui.resetKeybindings(); err != nil {
					return err
				}

				if err := gui.checkForChangedConfigsThatDontAutoReload(oldConfig, gui.Config.GetUserConfig()); err != nil {
					return err
				}
			}

			gui.c.Log.Info("Receiving focus - refreshing")
			refreshErr := gui.helpers.Refresh.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			if reloadErr != nil {
				// An error from reloading the config is the more important one
				// to report to the user
				return reloadErr
			}
			return refreshErr
		}

		return nil
	})

	gui.g.SetOpenHyperlinkFunc(func(url string) error {
		if strings.HasPrefix(url, "lazygit-edit:") {
			re := regexp.MustCompile(`^lazygit-edit://(.+?)(?::(\d+))?$`)
			matches := re.FindStringSubmatch(url)
			if matches == nil {
				return fmt.Errorf(gui.Tr.InvalidLazygitEditURL, url)
			}
			filepath := matches[1]
			if matches[2] != "" {
				lineNumber := utils.MustConvertToInt(matches[2])
				return gui.helpers.Files.EditFileAtLine(filepath, lineNumber)
			}
			return gui.helpers.Files.EditFiles([]string{filepath})
		}

		if err := gui.os.OpenLink(url); err != nil {
			return fmt.Errorf(gui.Tr.FailedToOpenURL, url, err)
		}

		return nil
	})

	// if a context key has been given, push that instead, and set its index to 0
	if contextKey != context.NO_CONTEXT {
		contextToPush = gui.c.ContextForKey(contextKey)
		// when we pass a list context, the expectation is that our cursor goes to the top,
		// because e.g. with worktrees, we'll show the current worktree at the top of the list.
		listContext, ok := contextToPush.(types.IListContext)
		if ok {
			listContext.GetList().SetSelection(0)
		}
	}

	gui.c.Context().Push(contextToPush)

	return nil
}

func (gui *Gui) getPerRepoConfigFiles() []*config.ConfigFile {
	repoConfigFiles := []*config.ConfigFile{
		// TODO: add filepath.Join(gui.git.RepoPaths.RepoPath(), ".lazygit.yml"),
		// with trust prompt
		{
			Path:   filepath.Join(gui.git.RepoPaths.RepoGitDirPath(), "lazygit.yml"),
			Policy: config.ConfigFilePolicySkipIfMissing,
		},
	}

	prevDir := gui.c.Git().RepoPaths.RepoPath()
	dir := filepath.Dir(prevDir)
	for dir != prevDir {
		repoConfigFiles = utils.Prepend(repoConfigFiles, &config.ConfigFile{
			Path:   filepath.Join(dir, ".lazygit.yml"),
			Policy: config.ConfigFilePolicySkipIfMissing,
		})
		prevDir = dir
		dir = filepath.Dir(dir)
	}
	return repoConfigFiles
}

func (gui *Gui) onUserConfigLoaded() error {
	userConfig := gui.Config.GetUserConfig()
	gui.Common.SetUserConfig(userConfig)

	gui.setColorScheme()
	gui.configureViewProperties()

	gui.g.SearchEscapeKey = keybindings.GetKey(userConfig.Keybinding.Universal.Return)
	gui.g.NextSearchMatchKey = keybindings.GetKey(userConfig.Keybinding.Universal.NextMatch)
	gui.g.PrevSearchMatchKey = keybindings.GetKey(userConfig.Keybinding.Universal.PrevMatch)

	gui.g.ShowListFooter = userConfig.Gui.ShowListFooter

	gui.g.Mouse = userConfig.Gui.MouseEvents

	if gui.previousLanguageConfig != userConfig.Gui.Language {
		tr, err := i18n.NewTranslationSetFromConfig(gui.Log, userConfig.Gui.Language)
		if err != nil {
			return err
		}
		gui.c.Tr = tr
		gui.previousLanguageConfig = userConfig.Gui.Language
	}

	// originally we could only hide the command log permanently via the config
	// but now we do it via state. So we need to still support the config for the
	// sake of backwards compatibility. We're making use of short circuiting here
	gui.ShowExtrasWindow = userConfig.Gui.ShowCommandLog && !gui.c.GetAppState().HideCommandLog

	authors.SetCustomAuthors(userConfig.Gui.AuthorColors)
	if userConfig.Gui.NerdFontsVersion != "" {
		icons.SetNerdFontsVersion(userConfig.Gui.NerdFontsVersion)
	} else if userConfig.Gui.ShowIcons {
		icons.SetNerdFontsVersion("2")
	}
	presentation.SetCustomBranches(userConfig.Gui.BranchColors)

	return nil
}

func (gui *Gui) checkForChangedConfigsThatDontAutoReload(oldConfig *config.UserConfig, newConfig *config.UserConfig) error {
	configsThatDontAutoReload := []string{
		"Git.AutoFetch",
		"Git.AutoRefresh",
		"Refresher.RefreshInterval",
		"Refresher.FetchInterval",
		"Update.Method",
		"Update.Days",
	}

	changedConfigs := []string{}
	for _, config := range configsThatDontAutoReload {
		old := reflect.ValueOf(oldConfig).Elem()
		new := reflect.ValueOf(newConfig).Elem()
		fieldNames := strings.Split(config, ".")
		userFacingPath := make([]string, 0, len(fieldNames))
		// navigate to the leaves in old and new config
		for _, fieldName := range fieldNames {
			f, _ := old.Type().FieldByName(fieldName)
			userFacingName := f.Tag.Get("yaml")
			if userFacingName == "" {
				userFacingName = fieldName
			}
			userFacingPath = append(userFacingPath, userFacingName)
			old = old.FieldByName(fieldName)
			new = new.FieldByName(fieldName)
		}
		// if the value has changed, ...
		if !old.Equal(new) {
			// ... append it to the list of changed configs
			changedConfigs = append(changedConfigs, strings.Join(userFacingPath, "."))
		}
	}

	if len(changedConfigs) == 0 {
		return nil
	}

	message := utils.ResolvePlaceholderString(
		gui.c.Tr.NonReloadableConfigWarning,
		map[string]string{
			"configs": strings.Join(changedConfigs, "\n"),
		},
	)
	gui.c.Confirm(types.ConfirmOpts{
		Title:  gui.c.Tr.NonReloadableConfigWarningTitle,
		Prompt: message,
	})

	return nil
}

// resetState reuses the repo state from our repo state map, if the repo was
// open before; otherwise it creates a new one.
func (gui *Gui) resetState(startArgs appTypes.StartArgs) types.Context {
	worktreePath := gui.git.RepoPaths.WorktreePath()

	if state := gui.RepoStateMap[Repo(worktreePath)]; state != nil {
		gui.State = state
		gui.State.ViewsSetup = false

		contextTree := gui.State.Contexts
		gui.State.WindowViewNameMap = initialWindowViewNameMap(contextTree)

		// setting this to nil so we don't get stuck based on a popup that was
		// previously opened
		gui.Mutexes.PopupMutex.Lock()
		gui.State.CurrentPopupOpts = nil
		gui.Mutexes.PopupMutex.Unlock()

		return gui.c.Context().Current()
	}

	contextTree := gui.contextTree()

	initialScreenMode := initialScreenMode(startArgs, gui.Config)

	gui.State = &GuiRepoState{
		ViewsSetup: false,
		Model: &types.Model{
			CommitFiles:           nil,
			Files:                 make([]*models.File, 0),
			Commits:               make([]*models.Commit, 0),
			StashEntries:          make([]*models.StashEntry, 0),
			FilteredReflogCommits: make([]*models.Commit, 0),
			ReflogCommits:         make([]*models.Commit, 0),
			BisectInfo:            git_commands.NewNullBisectInfo(),
			FilesTrie:             patricia.NewTrie(),
			Authors:               map[string]*models.Author{},
			MainBranches:          git_commands.NewMainBranches(gui.c.Common, gui.os.Cmd),
		},
		Modes: &types.Modes{
			Filtering:        filtering.New(startArgs.FilterPath, ""),
			CherryPicking:    cherrypicking.New(),
			Diffing:          diffing.New(),
			MarkedBaseCommit: marked_base_commit.New(),
		},
		ScreenMode: initialScreenMode,
		// TODO: only use contexts from context manager
		ContextMgr:        NewContextMgr(gui, contextTree),
		Contexts:          contextTree,
		WindowViewNameMap: initialWindowViewNameMap(contextTree),
		SearchState:       types.NewSearchState(),
	}

	gui.RepoStateMap[Repo(worktreePath)] = gui.State

	return initialContext(contextTree, startArgs)
}

func initialWindowViewNameMap(contextTree *context.ContextTree) *utils.ThreadSafeMap[string, string] {
	result := utils.NewThreadSafeMap[string, string]()

	for _, context := range contextTree.Flatten() {
		result.Set(context.GetWindowName(), context.GetViewName())
	}

	return result
}

func initialScreenMode(startArgs appTypes.StartArgs, config config.AppConfigurer) types.WindowMaximisation {
	if startArgs.FilterPath != "" || startArgs.GitArg != appTypes.GitArgNone {
		return types.SCREEN_FULL
	} else {
		defaultWindowSize := config.GetUserConfig().Gui.WindowSize

		switch defaultWindowSize {
		case "half":
			return types.SCREEN_HALF
		case "full":
			return types.SCREEN_FULL
		default:
			return types.SCREEN_NORMAL
		}
	}
}

func initialContext(contextTree *context.ContextTree, startArgs appTypes.StartArgs) types.IListContext {
	var initialContext types.IListContext = contextTree.Files

	if startArgs.FilterPath != "" {
		initialContext = contextTree.LocalCommits
	} else if startArgs.GitArg != appTypes.GitArgNone {
		switch startArgs.GitArg {
		case appTypes.GitArgStatus:
			initialContext = contextTree.Files
		case appTypes.GitArgBranch:
			initialContext = contextTree.Branches
		case appTypes.GitArgLog:
			initialContext = contextTree.LocalCommits
		case appTypes.GitArgStash:
			initialContext = contextTree.Stash
		default:
			panic("unhandled git arg")
		}
	}

	return initialContext
}

func (gui *Gui) Contexts() *context.ContextTree {
	return gui.State.Contexts
}

// for now the split view will always be on
// NewGui builds a new gui handler
func NewGui(
	cmn *common.Common,
	config config.AppConfigurer,
	gitVersion *git_commands.GitVersion,
	updater *updates.Updater,
	showRecentRepos bool,
	initialDir string,
	test integrationTypes.IntegrationTest,
) (*Gui, error) {
	gui := &Gui{
		Common:               cmn,
		gitVersion:           gitVersion,
		Config:               config,
		Updater:              updater,
		statusManager:        status.NewStatusManager(),
		viewBufferManagerMap: map[string]*tasks.ViewBufferManager{},
		viewPtmxMap:          map[string]*os.File{},
		showRecentRepos:      showRecentRepos,
		RepoPathStack:        &utils.StringStack{},
		RepoStateMap:         map[Repo]*GuiRepoState{},
		GuiLog:               []string{},

		// initializing this to true for the time being; it will be reset to the
		// real value after loading the user config:
		ShowExtrasWindow: true,

		Mutexes: types.Mutexes{
			RefreshingFilesMutex:    &deadlock.Mutex{},
			RefreshingBranchesMutex: &deadlock.Mutex{},
			RefreshingStatusMutex:   &deadlock.Mutex{},
			LocalCommitsMutex:       &deadlock.Mutex{},
			SubCommitsMutex:         &deadlock.Mutex{},
			AuthorsMutex:            &deadlock.Mutex{},
			SubprocessMutex:         &deadlock.Mutex{},
			PopupMutex:              &deadlock.Mutex{},
			PtyMutex:                &deadlock.Mutex{},
		},
		InitialDir:       initialDir,
		afterLayoutFuncs: make(chan func() error, 1000),

		itemOperations:      make(map[string]types.ItemOperation),
		itemOperationsMutex: &deadlock.Mutex{},
	}

	gui.PopupHandler = popup.NewPopupHandler(
		cmn,
		func(ctx goContext.Context, opts types.CreatePopupPanelOpts) {
			gui.helpers.Confirmation.CreatePopupPanel(ctx, opts)
		},
		func() error { return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC}) },
		func() { gui.State.ContextMgr.Pop() },
		func() types.Context { return gui.State.ContextMgr.Current() },
		gui.createMenu,
		func(message string, f func(gocui.Task) error) { gui.helpers.AppStatus.WithWaitingStatus(message, f) },
		func(message string, f func() error) error {
			return gui.helpers.AppStatus.WithWaitingStatusSync(message, f)
		},
		func(message string, kind types.ToastKind) { gui.helpers.AppStatus.Toast(message, kind) },
		func() string { return gui.Views.Confirmation.TextArea.GetContent() },
		func() bool { return gui.c.InDemo() },
	)

	guiCommon := &guiCommon{gui: gui, IPopupHandler: gui.PopupHandler}
	helperCommon := &helpers.HelperCommon{IGuiCommon: guiCommon, Common: cmn, IGetContexts: gui}

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

	gui.BackgroundRoutineMgr = &BackgroundRoutineMgr{gui: gui}
	gui.stateAccessor = &StateAccessor{gui: gui}

	return gui, nil
}

var RuneReplacements = map[rune]string{
	// for the commit graph
	graph.MergeSymbol:  "M",
	graph.CommitSymbol: "o",
}

func (gui *Gui) initGocui(headless bool, test integrationTypes.IntegrationTest) (*gocui.Gui, error) {
	runInSandbox := os.Getenv(components.SANDBOX_ENV_VAR) == "true"
	playRecording := test != nil && !runInSandbox

	width, height := 0, 0
	if test != nil {
		if test.RequiresHeadless() {
			if runInSandbox {
				panic("Test requires headless, can't run in sandbox")
			}
			headless = true
		}
		width, height = test.HeadlessDimensions()
	}

	g, err := gocui.NewGui(gocui.NewGuiOpts{
		OutputMode:       gocui.OutputTrue,
		SupportOverlaps:  OverlappingEdges,
		PlayRecording:    playRecording,
		Headless:         headless,
		RuneReplacements: RuneReplacements,
		Width:            width,
		Height:           height,
	})
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (gui *Gui) viewTabMap() map[string][]context.TabView {
	result := map[string][]context.TabView{
		"branches": {
			{
				Tab:      gui.c.Tr.LocalBranchesTitle,
				ViewName: "localBranches",
			},
			{
				Tab:      gui.c.Tr.RemotesTitle,
				ViewName: "remotes",
			},
			{
				Tab:      gui.c.Tr.TagsTitle,
				ViewName: "tags",
			},
		},
		"commits": {
			{
				Tab:      gui.c.Tr.CommitsTitle,
				ViewName: "commits",
			},
			{
				Tab:      gui.c.Tr.ReflogCommitsTitle,
				ViewName: "reflogCommits",
			},
		},
		"files": {
			{
				Tab:      gui.c.Tr.FilesTitle,
				ViewName: "files",
			},
			context.TabView{
				Tab:      gui.c.Tr.WorktreesTitle,
				ViewName: "worktrees",
			},
			{
				Tab:      gui.c.Tr.SubmodulesTitle,
				ViewName: "submodules",
			},
		},
	}

	return result
}

// Run: setup the gui with keybindings and start the mainloop
func (gui *Gui) Run(startArgs appTypes.StartArgs) error {
	g, err := gui.initGocui(Headless(), startArgs.IntegrationTest)
	if err != nil {
		return err
	}

	defer gui.checkForDeprecatedEditConfigs()

	gui.g = g
	defer gui.g.Close()

	g.ErrorHandler = gui.PopupHandler.ErrorHandler

	// if the deadlock package wants to report a deadlock, we first need to
	// close the gui so that we can actually read what it prints.
	deadlock.Opts.LogBuf = utils.NewOnceWriter(os.Stderr, func() {
		gui.g.Close()
	})
	// disable deadlock reporting if we're not running in debug mode, or if
	// we're debugging an integration test. In this latter case, stopping at
	// breakpoints and stepping through code can easily take more than 30s.
	deadlock.Opts.Disable = !gui.Debug || os.Getenv(components.WAIT_FOR_DEBUGGER_ENV_VAR) != ""

	gui.g.OnSearchEscape = func() error { gui.helpers.Search.Cancel(); return nil }

	gui.g.SetManager(gocui.ManagerFunc(gui.layout))

	if err := gui.createAllViews(); err != nil {
		return err
	}

	// onNewRepo must be called after g.SetManager because SetManager deletes keybindings
	if err := gui.onNewRepo(startArgs, context.NO_CONTEXT); err != nil {
		return err
	}

	gui.waitForIntro.Add(1)

	gui.BackgroundRoutineMgr.startBackgroundRoutines()

	gui.c.Log.Info("starting main loop")

	// setting here so we can use it in layout.go
	gui.integrationTest = startArgs.IntegrationTest

	return gui.g.MainLoop()
}

func (gui *Gui) RunAndHandleError(startArgs appTypes.StartArgs) error {
	gui.stopChan = make(chan struct{})
	return utils.SafeWithError(func() error {
		if err := gui.Run(startArgs); err != nil {
			for _, manager := range gui.viewBufferManagerMap {
				manager.Close()
			}

			close(gui.stopChan)

			switch err {
			case gocui.ErrQuit:
				if gui.c.State().GetRetainOriginalDir() {
					if err := gui.helpers.RecordDirectory.RecordDirectory(gui.InitialDir); err != nil {
						return err
					}
				} else {
					if err := gui.helpers.RecordDirectory.RecordCurrentDirectory(); err != nil {
						return err
					}
				}

				return nil

			default:
				return err
			}
		}

		return nil
	})
}

func (gui *Gui) checkForDeprecatedEditConfigs() {
	osConfig := &gui.UserConfig().OS
	deprecatedConfigs := []struct {
		config  string
		oldName string
		newName string
	}{
		{osConfig.EditCommand, "EditCommand", "Edit"},
		{osConfig.EditCommandTemplate, "EditCommandTemplate", "Edit,EditAtLine"},
		{osConfig.OpenCommand, "OpenCommand", "Open"},
		{osConfig.OpenLinkCommand, "OpenLinkCommand", "OpenLink"},
	}
	deprecatedConfigStrings := []string{}

	for _, dc := range deprecatedConfigs {
		if dc.config != "" {
			deprecatedConfigStrings = append(deprecatedConfigStrings, fmt.Sprintf("   OS.%s -> OS.%s", dc.oldName, dc.newName))
		}
	}
	if len(deprecatedConfigStrings) != 0 {
		warningMessage := utils.ResolvePlaceholderString(
			gui.c.Tr.DeprecatedEditConfigWarning,
			map[string]string{
				"configs": strings.Join(deprecatedConfigStrings, "\n"),
			},
		)

		os.Stdout.Write([]byte(warningMessage))
	}
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

	if err := gui.g.Suspend(); err != nil {
		return false, err
	}

	gui.BackgroundRoutineMgr.PauseBackgroundRefreshes(true)
	defer gui.BackgroundRoutineMgr.PauseBackgroundRefreshes(false)

	cmdErr := gui.runSubprocess(subprocess)

	if err := gui.g.Resume(); err != nil {
		return false, err
	}

	if cmdErr != nil {
		return false, cmdErr
	}

	return true, nil
}

func (gui *Gui) runSubprocess(cmdObj oscommands.ICmdObj) error { //nolint:unparam
	gui.LogCommand(cmdObj.ToString(), true)

	subprocess := cmdObj.GetCmd()
	subprocess.Stdout = os.Stdout
	subprocess.Stderr = os.Stderr
	subprocess.Stdin = os.Stdin

	fmt.Fprintf(os.Stdout, "\n%s\n\n", style.FgBlue.Sprint("+ "+strings.Join(subprocess.Args, " ")))

	err := subprocess.Run()

	subprocess.Stdout = io.Discard
	subprocess.Stderr = io.Discard
	subprocess.Stdin = nil

	if gui.integrationTest == nil && (gui.Config.GetUserConfig().PromptToReturnFromSubprocess || err != nil) {
		fmt.Fprintf(os.Stdout, "\n%s", style.FgGreen.Sprint(gui.Tr.PressEnterToReturn))

		// scan to buffer to prevent run unintentional operations when TUI resumes.
		var buffer string
		_, _ = fmt.Scanln(&buffer) // wait for enter press
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

func (gui *Gui) showIntroPopupMessage() {
	gui.waitForIntro.Add(1)

	gui.c.OnUIThread(func() error {
		onConfirm := func() error {
			gui.c.GetAppState().StartupPopupVersion = StartupPopupVersion
			err := gui.c.SaveAppState()
			gui.waitForIntro.Done()
			return err
		}

		gui.c.Confirm(types.ConfirmOpts{
			Title:         "",
			Prompt:        gui.c.Tr.IntroPopupMessage,
			HandleConfirm: onConfirm,
			HandleClose:   onConfirm,
		})

		return nil
	})
}

func (gui *Gui) showBreakingChangesMessage() {
	_, err := types.ParseVersionNumber(gui.Config.GetVersion())
	if err != nil {
		// We don't have a parseable version, so we'll assume it's a developer
		// build, or a build from HEAD with a version such as 0.40.0-g1234567;
		// in these cases we don't show release notes.
		return
	}

	last := &types.VersionNumber{}
	lastVersionStr := gui.c.GetAppState().LastVersion
	// If there's no saved last version, we show all release notes. This is for
	// people upgrading from a version before we started to save lastVersion.
	// First time new users won't see the release notes because we show them the
	// intro popup instead.
	if lastVersionStr != "" {
		last, err = types.ParseVersionNumber(lastVersionStr)
		if err != nil {
			// The last version was a developer build, so don't show release
			// notes in this case either.
			return
		}
	}

	// Now collect all release notes texts for versions newer than lastVersion.
	// We don't need to bother checking the current version here, because we
	// can't possibly have texts for versions newer than current.
	type versionAndText struct {
		version *types.VersionNumber
		text    string
	}
	texts := []versionAndText{}
	for versionStr, text := range gui.Tr.BreakingChangesByVersion {
		v, err := types.ParseVersionNumber(versionStr)
		if err != nil {
			// Ignore bogus entries in the BreakingChanges map
			continue
		}
		if last.IsOlderThan(v) {
			texts = append(texts, versionAndText{version: v, text: text})
		}
	}

	if len(texts) > 0 {
		sort.Slice(texts, func(i, j int) bool {
			return texts[i].version.IsOlderThan(texts[j].version)
		})
		message := strings.Join(lo.Map(texts, func(t versionAndText, _ int) string { return t.text }), "\n")

		gui.waitForIntro.Add(1)
		gui.c.OnUIThread(func() error {
			onConfirm := func() error {
				gui.waitForIntro.Done()
				return nil
			}

			gui.c.Confirm(types.ConfirmOpts{
				Title:         gui.Tr.BreakingChangesTitle,
				Prompt:        gui.Tr.BreakingChangesMessage + "\n\n" + message,
				HandleConfirm: onConfirm,
				HandleClose:   onConfirm,
			})
			return nil
		})
	}
}

// setColorScheme sets the color scheme for the app based on the user config
func (gui *Gui) setColorScheme() {
	userConfig := gui.UserConfig()
	theme.UpdateTheme(userConfig.Gui.Theme)

	gui.g.FgColor = theme.InactiveBorderColor
	gui.g.SelFgColor = theme.ActiveBorderColor
	gui.g.FrameColor = theme.InactiveBorderColor
	gui.g.SelFrameColor = theme.ActiveBorderColor
}

func (gui *Gui) onUIThread(f func() error) {
	gui.g.Update(func(*gocui.Gui) error {
		return f()
	})
}

func (gui *Gui) onWorker(f func(gocui.Task) error) {
	gui.g.OnWorker(f)
}

func (gui *Gui) getWindowDimensions(informationStr string, appStatus string) map[string]boxlayout.Dimensions {
	return gui.helpers.WindowArrangement.GetWindowDimensions(informationStr, appStatus)
}

func (gui *Gui) afterLayout(f func() error) {
	select {
	case gui.afterLayoutFuncs <- f:
	default:
		// hopefully this never happens
		gui.c.Log.Error("afterLayoutFuncs channel is full, skipping function")
	}
}
