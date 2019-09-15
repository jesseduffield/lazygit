package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shibukawa/configdir"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// AppConfig contains the base configuration fields required for lazygit.
type AppConfig struct {
	Debug         bool   `long:"debug" env:"DEBUG" default:"false"`
	Version       string `long:"version" env:"VERSION" default:"unversioned"`
	Commit        string `long:"commit" env:"COMMIT"`
	BuildDate     string `long:"build-date" env:"BUILD_DATE"`
	Name          string `long:"name" env:"NAME" default:"lazygit"`
	BuildSource   string `long:"build-source" env:"BUILD_SOURCE" default:""`
	UserConfig    *viper.Viper
	UserConfigDir string
	AppState      *AppState
	IsNewRepo     bool
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
	GetUserConfig() *viper.Viper
	GetUserConfigDir() string
	GetAppState() *AppState
	WriteToUserConfig(string, string) error
	SaveAppState() error
	LoadAppState() error
	SetIsNewRepo(bool)
	GetIsNewRepo() bool
}

// NewAppConfig makes a new app config
func NewAppConfig(name, version, commit, date string, buildSource string, debuggingFlag bool) (*AppConfig, error) {
	userConfig, userConfigPath, err := LoadConfig("config", true)
	if err != nil {
		return nil, err
	}

	if os.Getenv("DEBUG") == "TRUE" {
		debuggingFlag = true
	}

	appConfig := &AppConfig{
		Name:          "lazygit",
		Version:       version,
		Commit:        commit,
		BuildDate:     date,
		Debug:         debuggingFlag,
		BuildSource:   buildSource,
		UserConfig:    userConfig,
		UserConfigDir: filepath.Dir(userConfigPath),
		AppState:      &AppState{},
		IsNewRepo:     false,
	}

	if err := appConfig.LoadAppState(); err != nil {
		return nil, err
	}

	return appConfig, nil
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
func (c *AppConfig) GetUserConfig() *viper.Viper {
	return c.UserConfig
}

// GetAppState returns the app state
func (c *AppConfig) GetAppState() *AppState {
	return c.AppState
}

func (c *AppConfig) GetUserConfigDir() string {
	return c.UserConfigDir
}

func newViper(filename string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(filename)
	return v, nil
}

// LoadConfig gets the user's config
func LoadConfig(filename string, withDefaults bool) (*viper.Viper, string, error) {
	v, err := newViper(filename)
	if err != nil {
		return nil, "", err
	}
	if withDefaults {
		if err = LoadDefaults(v, GetDefaultConfig()); err != nil {
			return nil, "", err
		}
		if err = LoadDefaults(v, GetPlatformDefaultConfig()); err != nil {
			return nil, "", err
		}
	}
	configPath, err := LoadAndMergeFile(v, filename+".yml")
	if err != nil {
		return nil, "", err
	}
	return v, configPath, nil
}

// LoadDefaults loads in the defaults defined in this file
func LoadDefaults(v *viper.Viper, defaults []byte) error {
	return v.MergeConfig(bytes.NewBuffer(defaults))
}

func prepareConfigFile(filename string) (string, error) {
	// chucking my name there is not for vanity purposes, the xdg spec (and that
	// function) requires a vendor name. May as well line up with github
	configDirs := configdir.New("jesseduffield", "lazygit")
	folder := configDirs.QueryFolderContainsFile(filename)
	if folder == nil {
		// create the file as empty
		folders := configDirs.QueryFolders(configdir.Global)
		if err := folders[0].WriteFile(filename, []byte{}); err != nil {
			return "", err
		}
		folder = configDirs.QueryFolderContainsFile(filename)
	}
	return filepath.Join(folder.Path, filename), nil
}

// LoadAndMergeFile Loads the config/state file, creating
// the file has an empty one if it does not exist
func LoadAndMergeFile(v *viper.Viper, filename string) (string, error) {
	configPath, err := prepareConfigFile(filename)
	if err != nil {
		return "", err
	}

	v.AddConfigPath(filepath.Dir(configPath))
	return configPath, v.MergeInConfig()
}

// WriteToUserConfig adds a key/value pair to the user's config and saves it
func (c *AppConfig) WriteToUserConfig(key, value string) error {
	// reloading the user config directly (without defaults) so that we're not
	// writing any defaults back to the user's config
	v, _, err := LoadConfig("config", false)
	if err != nil {
		return err
	}

	v.Set(key, value)
	return v.WriteConfig()
}

// SaveAppState marshalls the AppState struct and writes it to the disk
func (c *AppConfig) SaveAppState() error {
	marshalledAppState, err := yaml.Marshal(c.AppState)
	if err != nil {
		return err
	}

	filepath, err := prepareConfigFile("state.yml")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, marshalledAppState, 0644)
}

// LoadAppState loads recorded AppState from file
func (c *AppConfig) LoadAppState() error {
	filepath, err := prepareConfigFile("state.yml")
	if err != nil {
		return err
	}
	appStateBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	if len(appStateBytes) == 0 {
		return yaml.Unmarshal(getDefaultAppState(), c.AppState)
	}
	return yaml.Unmarshal(appStateBytes, c.AppState)
}

// GetDefaultConfig returns the application default configuration
func GetDefaultConfig() []byte {
	return []byte(
		`gui:
  ## stuff relating to the UI
  scrollHeight: 2
  scrollPastBottom: true
  mouseEvents: false # will default to true when the feature is complete
  theme:
    activeBorderColor:
      - white
      - bold
    inactiveBorderColor:
      - white
    optionsTextColor:
      - blue
  commitLength:
    show: true
git:
  merging:
    manualCommit: false
  skipHookPrefix: 'WIP'
  autoFetch: true
update:
  method: prompt # can be: prompt | background | never
  days: 14 # how often a update is checked for
reporting: 'undetermined' # one of: 'on' | 'off' | 'undetermined'
confirmOnQuit: false
`)
}

// AppState stores data between runs of the app like when the last update check
// was performed and which other repos have been checked out
type AppState struct {
	LastUpdateCheck int64
	RecentRepos     []string
}

func getDefaultAppState() []byte {
	return []byte(`
    lastUpdateCheck: 0
    recentRepos: []
  `)
}

// // commenting this out until we use it again
// func homeDirectory() string {
// 	usr, err := user.Current()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return usr.HomeDir
// }
