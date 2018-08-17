package config

import (
	"bytes"
	"log"
	"os/user"

	"github.com/rollbar/rollbar-go"
	"github.com/shibukawa/configdir"
	"github.com/spf13/viper"
)

// AppConfig contains the base configuration fields required for lazygit.
type AppConfig struct {
	Debug      bool   `long:"debug" env:"DEBUG" default:"false"`
	Version    string `long:"version" env:"VERSION" default:"unversioned"`
	Commit     string `long:"commit" env:"COMMIT"`
	BuildDate  string `long:"build-date" env:"BUILD_DATE"`
	Name       string `long:"name" env:"NAME" default:"lazygit"`
	UserConfig *viper.Viper
}

// AppConfigurer interface allows individual app config structs to inherit Fields
// from AppConfig and still be used by lazygit.
type AppConfigurer interface {
	GetDebug() bool
	GetVersion() string
	GetCommit() string
	GetBuildDate() string
	GetName() string
	GetUserConfig() *viper.Viper
}

// NewAppConfig makes a new app config
func NewAppConfig(name, version, commit, date string, debuggingFlag *bool) (*AppConfig, error) {
	userConfig, err := LoadUserConfig()
	if err != nil {
		panic(err)
	}

	rollbar.SetToken("23432119147a4367abf7c0de2aa99a2d")
	rollbar.SetEnvironment("production") // defaults to "development"
	rollbar.SetCodeVersion(version)
	// rollbar.SetServerHost("web.1")                       // optional override; defaults to hostname
	rollbar.SetServerRoot("github.com/jesseduffield/lazygit") // path of project (required for GitHub integration and non-project stacktrace collapsing)

	// result, err := DoSomething()
	// if err != nil {
	//   rollbar.Critical(err)
	// }

	rollbar.Info("This is a test message")
	rollbar.Critical("test error")

	// rollbar.Wait()

	appConfig := &AppConfig{
		Name:       "lazygit",
		Version:    version,
		Commit:     commit,
		BuildDate:  date,
		Debug:      *debuggingFlag,
		UserConfig: userConfig,
	}
	return appConfig, nil
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

// GetUserConfig returns the user config
func (c *AppConfig) GetUserConfig() *viper.Viper {
	return c.UserConfig
}

// LoadUserConfig gets the user's config
func LoadUserConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	defaults := getDefaultConfig()
	err := v.ReadConfig(bytes.NewBuffer(defaults))
	if err != nil {
		return nil, err
	}
	v.SetConfigName("config")

	// chucking my name there is not for vanity purposes, the xdg spec (and that
	// function) requires a vendor name. May as well line up with github
	configDirs := configdir.New("jesseduffield", "lazygit")
	folder := configDirs.QueryFolderContainsFile("config.yml")
	if folder != nil {
		configData, err := folder.ReadFile("config.yml")
		if err != nil {
			return nil, err
		}
		if err = v.MergeConfig(bytes.NewReader(configData)); err != nil {
			return nil, err
		}
	}
	return v, nil
}

func getDefaultConfig() []byte {
	return []byte(`
  gui:
    ## stuff relating to the UI
    scrollHeight: 2
  git:
    # stuff relating to git
  os:
    # stuff relating to the OS

`)
}

func homeDirectory() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
