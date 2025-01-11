package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/jesseduffield/lazygit/pkg/utils/yaml_utils"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// AppConfig contains the base configuration fields required for lazygit.
type AppConfig struct {
	debug                 bool   `long:"debug" env:"DEBUG" default:"false"`
	version               string `long:"version" env:"VERSION" default:"unversioned"`
	buildDate             string `long:"build-date" env:"BUILD_DATE"`
	name                  string `long:"name" env:"NAME" default:"lazygit"`
	buildSource           string `long:"build-source" env:"BUILD_SOURCE" default:""`
	userConfig            *UserConfig
	globalUserConfigFiles []*ConfigFile
	userConfigFiles       []*ConfigFile
	userConfigDir         string
	tempDir               string
	appState              *AppState
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
	ReloadUserConfigForRepo(repoConfigFiles []*ConfigFile) error
	ReloadChangedUserConfigFiles() (error, bool)
	GetTempDir() string

	GetAppState() *AppState
	SaveAppState() error
}

type ConfigFilePolicy int

const (
	ConfigFilePolicyCreateIfMissing ConfigFilePolicy = iota
	ConfigFilePolicyErrorIfMissing
	ConfigFilePolicySkipIfMissing
)

type ConfigFile struct {
	Path    string
	Policy  ConfigFilePolicy
	modDate time.Time
	exists  bool
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

	var configFiles []*ConfigFile
	customConfigFiles := os.Getenv("LG_CONFIG_FILE")
	if customConfigFiles != "" {
		// Load user defined config files
		userConfigPaths := strings.Split(customConfigFiles, ",")
		configFiles = lo.Map(userConfigPaths, func(path string, _ int) *ConfigFile {
			return &ConfigFile{Path: path, Policy: ConfigFilePolicyErrorIfMissing}
		})
	} else {
		// Load default config files
		path := filepath.Join(configDir, ConfigFilename)
		configFile := &ConfigFile{Path: path, Policy: ConfigFilePolicyCreateIfMissing}
		configFiles = []*ConfigFile{configFile}
	}

	userConfig, err := loadUserConfigWithDefaults(configFiles)
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
		name:                  name,
		version:               version,
		buildDate:             date,
		debug:                 debuggingFlag,
		buildSource:           buildSource,
		userConfig:            userConfig,
		globalUserConfigFiles: configFiles,
		userConfigFiles:       configFiles,
		userConfigDir:         configDir,
		tempDir:               tempDir,
		appState:              appState,
	}

	return appConfig, nil
}

func ConfigDir() string {
	_, filePath := findConfigFile("config.yml")

	return filepath.Dir(filePath)
}

func findOrCreateConfigDir() (string, error) {
	folder := ConfigDir()
	return folder, os.MkdirAll(folder, 0o755)
}

func loadUserConfigWithDefaults(configFiles []*ConfigFile) (*UserConfig, error) {
	return loadUserConfig(configFiles, GetDefaultConfig())
}

func loadUserConfig(configFiles []*ConfigFile, base *UserConfig) (*UserConfig, error) {
	for _, configFile := range configFiles {
		path := configFile.Path
		statInfo, err := os.Stat(path)
		if err == nil {
			configFile.exists = true
			configFile.modDate = statInfo.ModTime()
		} else {
			if !os.IsNotExist(err) {
				return nil, err
			}

			switch configFile.Policy {
			case ConfigFilePolicyErrorIfMissing:
				return nil, err

			case ConfigFilePolicySkipIfMissing:
				configFile.exists = false
				continue

			case ConfigFilePolicyCreateIfMissing:
				file, err := os.Create(path)
				if err != nil {
					if os.IsPermission(err) {
						// apparently when people have read-only permissions they prefer us to fail silently
						continue
					}
					return nil, err
				}
				file.Close()

				configFile.exists = true
				statInfo, err := os.Stat(configFile.Path)
				if err != nil {
					return nil, err
				}
				configFile.modDate = statInfo.ModTime()
			}
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		content, err = migrateUserConfig(path, content)
		if err != nil {
			return nil, err
		}

		existingCustomCommands := base.CustomCommands

		if err := yaml.Unmarshal(content, base); err != nil {
			return nil, fmt.Errorf("The config at `%s` couldn't be parsed, please inspect it before opening up an issue.\n%w", path, err)
		}

		base.CustomCommands = append(base.CustomCommands, existingCustomCommands...)

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
	changedContent := content

	pathsToReplace := []struct {
		oldPath []string
		newName string
	}{
		{[]string{"gui", "skipUnstageLineWarning"}, "skipDiscardChangeWarning"},
		{[]string{"keybinding", "universal", "executeCustomCommand"}, "executeShellCommand"},
		{[]string{"gui", "windowSize"}, "screenMode"},
	}

	var err error
	for _, pathToReplace := range pathsToReplace {
		changedContent, err = yaml_utils.RenameYamlKey(changedContent, pathToReplace.oldPath, pathToReplace.newName)
		if err != nil {
			return nil, fmt.Errorf("Couldn't migrate config file at `%s` for key %s: %s", path, strings.Join(pathToReplace.oldPath, "."), err)
		}
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
	return c.debug
}

func (c *AppConfig) GetVersion() string {
	return c.version
}

func (c *AppConfig) GetName() string {
	return c.name
}

// GetBuildSource returns the source of the build. For builds from goreleaser
// this will be binaryBuild
func (c *AppConfig) GetBuildSource() string {
	return c.buildSource
}

// GetUserConfig returns the user config
func (c *AppConfig) GetUserConfig() *UserConfig {
	return c.userConfig
}

// GetAppState returns the app state
func (c *AppConfig) GetAppState() *AppState {
	return c.appState
}

func (c *AppConfig) GetUserConfigPaths() []string {
	return lo.FilterMap(c.userConfigFiles, func(f *ConfigFile, _ int) (string, bool) {
		return f.Path, f.exists
	})
}

func (c *AppConfig) GetUserConfigDir() string {
	return c.userConfigDir
}

func (c *AppConfig) ReloadUserConfigForRepo(repoConfigFiles []*ConfigFile) error {
	configFiles := append(c.globalUserConfigFiles, repoConfigFiles...)
	userConfig, err := loadUserConfigWithDefaults(configFiles)
	if err != nil {
		return err
	}

	c.userConfig = userConfig
	c.userConfigFiles = configFiles
	return nil
}

func (c *AppConfig) ReloadChangedUserConfigFiles() (error, bool) {
	fileHasChanged := func(f *ConfigFile) bool {
		info, err := os.Stat(f.Path)
		if err != nil && !os.IsNotExist(err) {
			// If we can't stat the file, assume it hasn't changed
			return false
		}
		exists := err == nil
		return exists != f.exists || (exists && info.ModTime() != f.modDate)
	}

	if lo.NoneBy(c.userConfigFiles, fileHasChanged) {
		return nil, false
	}

	userConfig, err := loadUserConfigWithDefaults(c.userConfigFiles)
	if err != nil {
		return err, false
	}

	c.userConfig = userConfig
	return nil, true
}

func (c *AppConfig) GetTempDir() string {
	return c.tempDir
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

// SaveAppState marshalls the AppState struct and writes it to the disk
func (c *AppConfig) SaveAppState() error {
	marshalledAppState, err := yaml.Marshal(c.appState)
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

// SaveGlobalUserConfig saves the UserConfig back to disk. This is only used in
// integration tests, so we are a bit sloppy with error handling.
func (c *AppConfig) SaveGlobalUserConfig() {
	if len(c.globalUserConfigFiles) != 1 {
		panic("expected exactly one global user config file")
	}

	yamlContent, err := yaml.Marshal(c.userConfig)
	if err != nil {
		log.Fatalf("error marshalling user config: %v", err)
	}

	err = os.WriteFile(c.globalUserConfigFiles[0].Path, yamlContent, 0o644)
	if err != nil {
		log.Fatalf("error saving user config: %v", err)
	}
}

// AppState stores data between runs of the app like when the last update check
// was performed and which other repos have been checked out
type AppState struct {
	LastUpdateCheck     int64
	RecentRepos         []string
	StartupPopupVersion int
	LastVersion         string // this is the last version the user was using, for the purpose of showing release notes

	// these are for shell commands typed in directly, not for custom commands in the lazygit config.
	// For backwards compatibility we keep the old name in yaml files.
	ShellCommandsHistory []string `yaml:"customcommandshistory"`

	HideCommandLog             bool
	IgnoreWhitespaceInDiffView bool
	DiffContextSize            uint64
	RenameSimilarityThreshold  int
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
		LastUpdateCheck:           0,
		RecentRepos:               []string{},
		StartupPopupVersion:       0,
		LastVersion:               "",
		DiffContextSize:           3,
		RenameSimilarityThreshold: 50,
		LocalBranchSortOrder:      "recency",
		RemoteBranchSortOrder:     "alphabetical",
		GitLogOrder:               "", // should be "topo-order" eventually
		GitLogShowGraph:           "", // should be "always" eventually
	}
}

func LogPath() (string, error) {
	if os.Getenv("LAZYGIT_LOG_PATH") != "" {
		return os.Getenv("LAZYGIT_LOG_PATH"), nil
	}

	return stateFilePath("development.log")
}
