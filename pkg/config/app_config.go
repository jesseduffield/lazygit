package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/jesseduffield/lazygit/pkg/utils/yaml_utils"
	"gopkg.in/yaml.v3"
)

// AppConfig contains the base configuration fields required for lazygit.
type AppConfig struct {
	Debug            bool   `long:"debug" env:"DEBUG" default:"false"`
	Version          string `long:"version" env:"VERSION" default:"unversioned"`
	BuildDate        string `long:"build-date" env:"BUILD_DATE"`
	Name             string `long:"name" env:"NAME" default:"lazygit"`
	BuildSource      string `long:"build-source" env:"BUILD_SOURCE" default:""`
	UserConfig       *UserConfig
	UserConfigPaths  []string
	DeafultConfFiles bool
	UserConfigDir    string
	TempDir          string
	AppState         *AppState
	IsNewRepo        bool
}

type AppConfigurer interface {
	GetDebug() bool

	// build info
	GetVersion() string
	GetName() string
	GetBuildSource() string

	GetUserConfig() *UserConfig
	GetUserConfigPaths() []string
	GetUserConfigDir() string
	ReloadUserConfig() error
	GetTempDir() string

	GetAppState() *AppState
	SaveAppState() error
}

// NewAppConfig makes a new app config
func NewAppConfig(
	name string,
	version,
	commit,
	date string,
	buildSource string,
	debuggingFlag bool,
	tempDir string,
) (*AppConfig, error) {
	configDir, err := findOrCreateConfigDir()
	if err != nil && !os.IsPermission(err) {
		return nil, err
	}

	var userConfigPaths []string
	customConfigFiles := os.Getenv("LG_CONFIG_FILE")
	if customConfigFiles != "" {
		// Load user defined config files
		userConfigPaths = strings.Split(customConfigFiles, ",")
	} else {
		// Load default config files
		userConfigPaths = []string{filepath.Join(configDir, ConfigFilename)}
	}

	userConfig, err := loadUserConfigWithDefaults(userConfigPaths)
	if err != nil {
		return nil, err
	}

	appState, err := loadAppState()
	if err != nil {
		return nil, err
	}

	// Temporary: the defaults for these are set to empty strings in
	// getDefaultAppState so that we can migrate them from userConfig (which is
	// now deprecated). Once we remove the user configs, we can remove this code
	// and set the proper defaults in getDefaultAppState.
	if appState.GitLogOrder == "" {
		appState.GitLogOrder = userConfig.Git.Log.Order
	}
	if appState.GitLogShowGraph == "" {
		appState.GitLogShowGraph = userConfig.Git.Log.ShowGraph
	}

	appConfig := &AppConfig{
		Name:            name,
		Version:         version,
		BuildDate:       date,
		Debug:           debuggingFlag,
		BuildSource:     buildSource,
		UserConfig:      userConfig,
		UserConfigPaths: userConfigPaths,
		UserConfigDir:   configDir,
		TempDir:         tempDir,
		AppState:        appState,
		IsNewRepo:       false,
	}

	return appConfig, nil
}

func isCustomConfigFile(path string) bool {
	return path != filepath.Join(ConfigDir(), ConfigFilename)
}

func ConfigDir() string {
	_, filePath := findConfigFile("config.yml")

	return filepath.Dir(filePath)
}

func findOrCreateConfigDir() (string, error) {
	folder := ConfigDir()
	return folder, os.MkdirAll(folder, 0o755)
}

func loadUserConfigWithDefaults(configFiles []string) (*UserConfig, error) {
	return loadUserConfig(configFiles, GetDefaultConfig())
}

func loadUserConfig(configFiles []string, base *UserConfig) (*UserConfig, error) {
	for _, path := range configFiles {
		if _, err := os.Stat(path); err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}

			// if use has supplied their own custom config file path(s), we assume
			// the files have already been created, so we won't go and create them here.
			if isCustomConfigFile(path) {
				return nil, err
			}

			file, err := os.Create(path)
			if err != nil {
				if os.IsPermission(err) {
					// apparently when people have read-only permissions they prefer us to fail silently
					continue
				}
				return nil, err
			}
			file.Close()
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		content, err = migrateUserConfig(path, content)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(content, base); err != nil {
			return nil, fmt.Errorf("The config at `%s` couldn't be parsed, please inspect it before opening up an issue.\n%w", path, err)
		}

		if err := base.Validate(); err != nil {
			return nil, fmt.Errorf("The config at `%s` has a validation error.\n%w", path, err)
		}
	}

	return base, nil
}

// Do any backward-compatibility migrations of things that have changed in the
// config over time; examples are renaming a key to a better name, moving a key
// from one container to another, or changing the type of a key (e.g. from bool
// to an enum).
func migrateUserConfig(path string, content []byte) ([]byte, error) {
	changedContent, err := yaml_utils.RenameYamlKey(content, []string{"gui", "skipUnstageLineWarning"},
		"skipDiscardChangeWarning")
	if err != nil {
		return nil, fmt.Errorf("Couldn't migrate config file at `%s`: %s", path, err)
	}

	changedContent, err = changeNullKeybindingsToDisabled(changedContent)
	if err != nil {
		return nil, fmt.Errorf("Couldn't migrate config file at `%s`: %s", path, err)
	}

	// Add more migrations here...

	// Write config back if changed
	if string(changedContent) != string(content) {
		if err := os.WriteFile(path, changedContent, 0o644); err != nil {
			return nil, fmt.Errorf("Couldn't write migrated config back to `%s`: %s", path, err)
		}
		return changedContent, nil
	}

	return content, nil
}

func changeNullKeybindingsToDisabled(changedContent []byte) ([]byte, error) {
	return yaml_utils.Walk(changedContent, func(node *yaml.Node, path string) bool {
		if strings.HasPrefix(path, "keybinding.") && node.Kind == yaml.ScalarNode && node.Tag == "!!null" {
			node.Value = "<disabled>"
			node.Tag = "!!str"
			return true
		}
		return false
	})
}

func (c *AppConfig) GetDebug() bool {
	return c.Debug
}

func (c *AppConfig) GetVersion() string {
	return c.Version
}

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

// GetAppState returns the app state
func (c *AppConfig) GetAppState() *AppState {
	return c.AppState
}

func (c *AppConfig) GetUserConfigPaths() []string {
	return c.UserConfigPaths
}

func (c *AppConfig) GetUserConfigDir() string {
	return c.UserConfigDir
}

func (c *AppConfig) ReloadUserConfig() error {
	userConfig, err := loadUserConfigWithDefaults(c.UserConfigPaths)
	if err != nil {
		return err
	}

	c.UserConfig = userConfig
	return nil
}

func (c *AppConfig) GetTempDir() string {
	return c.TempDir
}

// findConfigFile looks for a possibly existing config file.
// This function does NOT create any folders or files.
func findConfigFile(filename string) (exists bool, path string) {
	if envConfigDir := os.Getenv("CONFIG_DIR"); envConfigDir != "" {
		return true, filepath.Join(envConfigDir, filename)
	}

	// look for jesseduffield/lazygit/filename in XDG_CONFIG_HOME and XDG_CONFIG_DIRS
	legacyConfigPath, err := xdg.SearchConfigFile(filepath.Join("jesseduffield", "lazygit", filename))
	if err == nil {
		return true, legacyConfigPath
	}

	// look for lazygit/filename in XDG_CONFIG_HOME and XDG_CONFIG_DIRS
	configFilepath, err := xdg.SearchConfigFile(filepath.Join("lazygit", filename))
	if err == nil {
		return true, configFilepath
	}

	return false, filepath.Join(xdg.ConfigHome, "lazygit", filename)
}

var ConfigFilename = "config.yml"

// stateFilePath looks for a possibly existing state file.
// if none exist, the default path is returned and all parent directories are created.
func stateFilePath(filename string) (string, error) {
	exists, legacyStateFile := findConfigFile(filename)
	if exists {
		return legacyStateFile, nil
	}

	// looks for XDG_STATE_HOME/lazygit/filename
	return xdg.StateFile(filepath.Join("lazygit", filename))
}

// ConfigFilename returns the filename of the default config file
func (c *AppConfig) ConfigFilename() string {
	return filepath.Join(c.UserConfigDir, ConfigFilename)
}

// SaveAppState marshalls the AppState struct and writes it to the disk
func (c *AppConfig) SaveAppState() error {
	marshalledAppState, err := yaml.Marshal(c.AppState)
	if err != nil {
		return err
	}

	filepath, err := stateFilePath(stateFileName)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, marshalledAppState, 0o644)
	if err != nil && os.IsPermission(err) {
		// apparently when people have read-only permissions they prefer us to fail silently
		return nil
	}

	return err
}

var stateFileName = "state.yml"

// loadAppState loads recorded AppState from file
func loadAppState() (*AppState, error) {
	appState := getDefaultAppState()

	filepath, err := stateFilePath(stateFileName)
	if err != nil {
		if os.IsPermission(err) {
			// apparently when people have read-only permissions they prefer us to fail silently
			return appState, nil
		}
		return nil, err
	}

	appStateBytes, err := os.ReadFile(filepath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if len(appStateBytes) == 0 {
		return appState, nil
	}

	err = yaml.Unmarshal(appStateBytes, appState)
	if err != nil {
		return nil, err
	}

	return appState, nil
}

// AppState stores data between runs of the app like when the last update check
// was performed and which other repos have been checked out
type AppState struct {
	LastUpdateCheck     int64
	RecentRepos         []string
	StartupPopupVersion int
	LastVersion         string // this is the last version the user was using, for the purpose of showing release notes

	// these are for custom commands typed in directly, not for custom commands in the lazygit config
	CustomCommandsHistory      []string
	HideCommandLog             bool
	IgnoreWhitespaceInDiffView bool
	DiffContextSize            int
	LocalBranchSortOrder       string
	RemoteBranchSortOrder      string

	// One of: 'date-order' | 'author-date-order' | 'topo-order' | 'default'
	// 'topo-order' makes it easier to read the git log graph, but commits may not
	// appear chronologically. See https://git-scm.com/docs/
	GitLogOrder string

	// This determines whether the git graph is rendered in the commits panel
	// One of 'always' | 'never' | 'when-maximised'
	GitLogShowGraph string
}

func getDefaultAppState() *AppState {
	return &AppState{
		LastUpdateCheck:       0,
		RecentRepos:           []string{},
		StartupPopupVersion:   0,
		LastVersion:           "",
		DiffContextSize:       3,
		LocalBranchSortOrder:  "recency",
		RemoteBranchSortOrder: "alphabetical",
		GitLogOrder:           "", // should be "topo-order" eventually
		GitLogShowGraph:       "", // should be "always" eventually
	}
}

func LogPath() (string, error) {
	if os.Getenv("LAZYGIT_LOG_PATH") != "" {
		return os.Getenv("LAZYGIT_LOG_PATH"), nil
	}

	return stateFilePath("development.log")
}
