package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
	yaml "github.com/jesseduffield/yaml"
)

// AppConfig contains the base configuration fields required for lazygit.
type AppConfig struct {
	Debug          bool   `long:"debug" env:"DEBUG" default:"false"`
	Version        string `long:"version" env:"VERSION" default:"unversioned"`
	Commit         string `long:"commit" env:"COMMIT"`
	BuildDate      string `long:"build-date" env:"BUILD_DATE"`
	Name           string `long:"name" env:"NAME" default:"lazygit"`
	BuildSource    string `long:"build-source" env:"BUILD_SOURCE" default:""`
	UserConfig     *UserConfig
	UserConfigDir  string
	UserConfigPath string
	AppState       *AppState
	IsNewRepo      bool
}

// AppConfigurer interface allows individual app config structs to inherit Fields
// from AppConfig and still be used by lazygit.
type AppConfigurer interface {
	GetDebug() bool
	GetVersion() string
	GetCommit() string
	GetBuildDate() string
	GetName() string
	GetBuildSource() string
	GetUserConfig() *UserConfig
	GetUserConfigDir() string
	GetUserConfigPath() string
	GetAppState() *AppState
	SaveAppState() error
	SetIsNewRepo(bool)
	GetIsNewRepo() bool
}

// NewAppConfig makes a new app config
func NewAppConfig(name, version, commit, date string, buildSource string, debuggingFlag bool) (*AppConfig, error) {
	configDir, err := findOrCreateConfigDir(name)
	if err != nil {
		return nil, err
	}

	userConfig, err := loadUserConfigWithDefaults(configDir)
	if err != nil {
		return nil, err
	}

	if os.Getenv("DEBUG") == "TRUE" {
		debuggingFlag = true
	}

	appState, err := loadAppState()
	if err != nil {
		return nil, err
	}

	appConfig := &AppConfig{
		Name:           "lazygit",
		Version:        version,
		Commit:         commit,
		BuildDate:      date,
		Debug:          debuggingFlag,
		BuildSource:    buildSource,
		UserConfig:     userConfig,
		UserConfigDir:  configDir,
		UserConfigPath: filepath.Join(configDir, "config.yml"),
		AppState:       appState,
		IsNewRepo:      false,
	}

	return appConfig, nil
}

func findOrCreateConfigDir(projectName string) (string, error) {
	// chucking my name there is not for vanity purposes, the xdg spec (and that
	// function) requires a vendor name. May as well line up with github
	configDirs := xdg.New("jesseduffield", projectName)
	folder := configDirs.ConfigHome()

	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return "", err
	}

	return folder, nil
}

func loadUserConfigWithDefaults(configDir string) (*UserConfig, error) {
	return loadUserConfig(configDir, GetDefaultConfig())
}

func loadUserConfig(configDir string, base *UserConfig) (*UserConfig, error) {
	fileName := filepath.Join(configDir, "config.yml")

	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(fileName)
			if err != nil {
				return nil, err
			}
			file.Close()
		} else {
			return nil, err
		}
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(content, base); err != nil {
		return nil, err
	}

	return base, nil
}

// GetIsNewRepo returns known repo boolean
func (c *AppConfig) GetIsNewRepo() bool {
	return c.IsNewRepo
}

// SetIsNewRepo set if the current repo is known
func (c *AppConfig) SetIsNewRepo(toSet bool) {
	c.IsNewRepo = toSet
}

// GetDebug returns debug flag
func (c *AppConfig) GetDebug() bool {
	return c.Debug
}

// GetVersion returns debug flag
func (c *AppConfig) GetVersion() string {
	return c.Version
}

// GetCommit returns debug flag
func (c *AppConfig) GetCommit() string {
	return c.Commit
}

// GetBuildDate returns debug flag
func (c *AppConfig) GetBuildDate() string {
	return c.BuildDate
}

// GetName returns debug flag
func (c *AppConfig) GetName() string {
	return c.Name
}

// GetBuildSource returns the source of the build. For builds from goreleaser
// this will be binaryBuild
func (c *AppConfig) GetBuildSource() string {
	return c.BuildSource
}

// GetUserConfig returns the user config
func (c *AppConfig) GetUserConfig() *UserConfig {
	return c.UserConfig
}

// GetUserConfig returns the user config
func (c *AppConfig) GetUserConfigPath() string {
	return c.UserConfigPath
}

// GetAppState returns the app state
func (c *AppConfig) GetAppState() *AppState {
	return c.AppState
}

func (c *AppConfig) GetUserConfigDir() string {
	return c.UserConfigDir
}

func configFilePath(filename string) (string, error) {
	folder, err := findOrCreateConfigDir("lazygit")
	if err != nil {
		return "", err
	}

	return filepath.Join(folder, filename), nil
}

// ConfigFilename returns the filename of the current config file
func (c *AppConfig) ConfigFilename() string {
	return filepath.Join(c.UserConfigDir, "config.yml")
}

// SaveAppState marshalls the AppState struct and writes it to the disk
func (c *AppConfig) SaveAppState() error {
	marshalledAppState, err := yaml.Marshal(c.AppState)
	if err != nil {
		return err
	}

	filepath, err := configFilePath("state.yml")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, marshalledAppState, 0644)
}

// loadAppState loads recorded AppState from file
func loadAppState() (*AppState, error) {
	filepath, err := configFilePath("state.yml")
	if err != nil {
		return nil, err
	}

	appStateBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	if len(appStateBytes) == 0 {
		return getDefaultAppState(), nil
	}

	appState := &AppState{}
	err = yaml.Unmarshal(appStateBytes, appState)
	if err != nil {
		return nil, err
	}

	return appState, nil
}

func GetDefaultConfig() *UserConfig {
	userConfig := &UserConfig{}

	if err := yaml.Unmarshal(GetDefaultConfigBytes(), userConfig); err != nil {
		panic(err)
	}

	userConfig.OS = GetPlatformDefaultConfig()

	return userConfig
}

// GetDefaultConfigBytes returns the application default configuration
func GetDefaultConfigBytes() []byte {
	return []byte(
		`gui:
  ## stuff relating to the UI
  scrollHeight: 2
  scrollPastBottom: true
  mouseEvents: true
  skipUnstageLineWarning: false
  skipStashWarning: true
  sidePanelWidth: 0.3333
  expandFocusedSidePanel: false
  mainPanelSplitMode: 'flexible' # one of 'horizontal' | 'flexible' | 'vertical'
  theme:
    lightTheme: false
    activeBorderColor:
      - green
      - bold
    inactiveBorderColor:
      - white
    optionsTextColor:
      - blue
    selectedLineBgColor:
      - default
    selectedRangeBgColor:
      - blue
  commitLength:
    show: true
git:
  paging:
    colorArg: always
    useConfig: false
  merging:
    manualCommit: false
    args: ""
  pull:
    mode: 'merge' # one of 'merge' | 'rebase' | 'ff-only'
  skipHookPrefix: 'WIP'
  autoFetch: true
  branchLogCmd: "git log --graph --color=always --abbrev-commit --decorate --date=relative --pretty=medium {{branchName}} --"
  overrideGpg: false # prevents lazygit from spawning a separate process when using GPG
  disableForcePushing: false
update:
  method: prompt # can be: prompt | background | never
  days: 14 # how often a update is checked for
reporting: 'undetermined' # one of: 'on' | 'off' | 'undetermined'
splashUpdatesIndex: 0
confirmOnQuit: false
quitOnTopLevelReturn: true
disableStartupPopups: false
keybinding:
  universal:
    quit: 'q'
    quit-alt1: '<c-c>'
    return: '<esc>'
    quitWithoutChangingDirectory: 'Q'
    togglePanel: '<tab>'
    prevItem: '<up>'
    nextItem: '<down>'
    prevItem-alt: 'k'
    nextItem-alt: 'j'
    prevPage: ','
    nextPage: '.'
    gotoTop: '<'
    gotoBottom: '>'
    prevBlock: '<left>'
    nextBlock: '<right>'
    prevBlock-alt: 'h'
    nextBlock-alt: 'l'
    nextMatch: 'n'
    prevMatch: 'N'
    startSearch: '/'
    optionMenu: 'x'
    optionMenu-alt1: '?'
    select: '<space>'
    goInto: '<enter>'
    confirm: '<enter>'
    confirm-alt1: 'y'
    remove: 'd'
    new: 'n'
    edit: 'e'
    openFile: 'o'
    scrollUpMain: '<pgup>'
    scrollDownMain: '<pgdown>'
    scrollUpMain-alt1: 'K'
    scrollDownMain-alt1: 'J'
    scrollUpMain-alt2: '<c-u>'
    scrollDownMain-alt2: '<c-d>'
    executeCustomCommand: ':'
    createRebaseOptionsMenu: 'm'
    pushFiles: 'P'
    pullFiles: 'p'
    refresh: 'R'
    createPatchOptionsMenu: '<c-p>'
    nextTab: ']'
    prevTab: '['
    nextScreenMode: '+'
    prevScreenMode: '_'
    undo: 'z'
    redo: '<c-z>'
    filteringMenu: <c-s>
    diffingMenu: 'W'
    diffingMenu-alt: '<c-e>'
    copyToClipboard: '<c-o>'
  status:
    checkForUpdate: 'u'
    recentRepos: '<enter>'
  files:
    commitChanges: 'c'
    commitChangesWithoutHook: 'w'
    amendLastCommit: 'A'
    commitChangesWithEditor: 'C'
    ignoreFile: 'i'
    refreshFiles: 'r'
    stashAllChanges: 's'
    viewStashOptions: 'S'
    toggleStagedAll: 'a'
    viewResetOptions: 'D'
    fetch: 'f'
  branches:
    createPullRequest: 'o'
    checkoutBranchByName: 'c'
    forceCheckoutBranch: 'F'
    rebaseBranch: 'r'
    renameBranch: 'R'
    mergeIntoCurrentBranch: 'M'
    viewGitFlowOptions: 'i'
    fastForward: 'f'
    pushTag: 'P'
    setUpstream: 'u'
    fetchRemote: 'f'
  commits:
    squashDown: 's'
    renameCommit: 'r'
    renameCommitWithEditor: 'R'
    viewResetOptions: 'g'
    markCommitAsFixup: 'f'
    createFixupCommit: 'F'
    squashAboveCommits: 'S'
    moveDownCommit: '<c-j>'
    moveUpCommit: '<c-k>'
    amendToCommit: 'A'
    pickCommit: 'p'
    revertCommit: 't'
    cherryPickCopy: 'c'
    cherryPickCopyRange: 'C'
    pasteCommits: 'v'
    tagCommit: 'T'
    checkoutCommit: '<space>'
    resetCherryPick: '<c-R>'
  stash:
    popStash: 'g'
  commitFiles:
    checkoutCommitFile: 'c'
  main:
    toggleDragSelect: 'v'
    toggleDragSelect-alt: 'V'
    toggleSelectHunk: 'a'
    pickBothHunks: 'b'
  submodules:
    init: 'i'
    update: 'u'
    bulkMenu: 'b'
`)
}

// AppState stores data between runs of the app like when the last update check
// was performed and which other repos have been checked out
type AppState struct {
	LastUpdateCheck     int64
	RecentRepos         []string
	StartupPopupVersion int
}

func getDefaultAppState() *AppState {
	return &AppState{
		LastUpdateCheck:     0,
		RecentRepos:         []string{},
		StartupPopupVersion: 0,
	}
}

func LogPath() (string, error) {
	return configFilePath("development.log")
}
