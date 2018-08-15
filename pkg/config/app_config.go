package config

import (
	"bytes"
	"log"
	"os"
	"os/user"
	"path/filepath"

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
	v.SetConfigType("json")
	defaults := getDefaultConfig()
	err := v.ReadConfig(bytes.NewBuffer(defaults))
	if err != nil {
		return nil, err
	}
	v.SetConfigName("config")
	configPath := homeDirectory() + "/lazygit/"
	if _, err := os.Stat(filepath.FromSlash(configPath + "config.json")); !os.IsNotExist(err) {
		v.AddConfigPath(configPath)
		err = v.MergeInConfig()
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func getDefaultConfig() []byte {
	return []byte(`
	{
		"gui": {
			"scrollHeight": 1
		},
		"git": {},
		"os": {}
	}
`)
}

func homeDirectory() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
