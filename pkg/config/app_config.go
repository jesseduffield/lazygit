package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenPeeDeeP/xdg"
	yaml "github.com/jesseduffield/yaml"
)

// AppConfig contains the base configuration fields required for lazygit.
type AppConfig struct {
	Debug            bool   `long:"debug" env:"DEBUG" default:"false"`
	Version          string `long:"version" env:"VERSION" default:"unversioned"`
	Commit           string `long:"commit" env:"COMMIT"`
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

// AppConfigurer interface allows individual app config structs to inherit Fields
// from AppConfig and still be used by lazygit.
type AppConfigurer interface {
	GetDebug() bool

	// build info
	GetVersion() string
	GetCommit() string
	GetBuildDate() string
	GetName() string
	GetBuildSource() string

	GetUserConfig() *UserConfig
	GetUserConfigPaths() []string
	GetUserConfigDir() string
	ReloadUserConfig() error

	GetAppState() *AppState
	SaveAppState() error
}

// NewAppConfig makes a new app config
func NewAppConfig(name, version, commit, date string, buildSource string, debuggingFlag bool) (*AppConfig, error) {
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

	if os.Getenv("DEBUG") == "TRUE" {
		debuggingFlag = true
	}

	tempDir := filepath.Join(os.TempDir(), "lazygit")

	appState, err := loadAppState()
	if err != nil {
		return nil, err
	}

	appConfig := &AppConfig{
		Name:            "lazygit",
		Version:         version,
		Commit:          commit,
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
	legacyConfigDirectory := configDirForVendor("jesseduffield")
	if _, err := os.Stat(legacyConfigDirectory); !os.IsNotExist(err) {
		return legacyConfigDirectory
	}
	configDirectory := configDirForVendor("")
	return configDirectory
}

func configDirForVendor(vendor string) string {
	envConfigDir := os.Getenv("CONFIG_DIR")
	if envConfigDir != "" {
		return envConfigDir
	}
	configDirs := xdg.New(vendor, "lazygit")
	return configDirs.ConfigHome()
}

func findOrCreateConfigDir() (string, error) {
	folder := ConfigDir()
	return folder, os.MkdirAll(folder, 0755)
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

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(content, base); err != nil {
			return nil, err
		}
	}

	return base, nil
}

func (c *AppConfig) GetDebug() bool {
	return c.Debug
}

func (c *AppConfig) GetVersion() string {
	return c.Version
}

func (c *AppConfig) GetCommit() string {
	return c.Commit
}

func (c *AppConfig) GetBuildDate() string {
	return c.BuildDate
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

func configFilePath(filename string) (string, error) {
	folder, err := findOrCreateConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(folder, filename), nil
}

var ConfigFilename = "config.yml"

// ConfigFilename returns the filename of the deafult config file
func (c *AppConfig) ConfigFilename() string {
	return filepath.Join(c.UserConfigDir, ConfigFilename)
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

	err = ioutil.WriteFile(filepath, marshalledAppState, 0644)
	if err != nil && os.IsPermission(err) {
		// apparently when people have read-only permissions they prefer us to fail silently
		return nil
	}

	return err
}

// loadAppState loads recorded AppState from file
func loadAppState() (*AppState, error) {
	filepath, err := configFilePath("state.yml")
	if err != nil {
		if os.IsPermission(err) {
			// apparently when people have read-only permissions they prefer us to fail silently
			return getDefaultAppState(), nil
		}
		return nil, err
	}

	appStateBytes, err := ioutil.ReadFile(filepath)
	if err != nil && !os.IsNotExist(err) {
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

// AppState stores data between runs of the app like when the last update check
// was performed and which other repos have been checked out
type AppState struct {
	LastUpdateCheck     int64
	RecentRepos         []string
	StartupPopupVersion int

	// these are for custom commands typed in directly, not for custom commands in the lazygit config
	CustomCommandsHistory []string
	HideCommandLog        bool
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
