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
	configDir, err := findOrCreateConfigDir()
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
				if strings.Contains(err.Error(), "read-only file system") {
					return base, nil
				}
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
	folder, err := findOrCreateConfigDir()
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
