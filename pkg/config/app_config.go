package config

// AppConfig contains the base configuration fields required for lazygit.
type AppConfig struct {
	Debug     bool   `long:"debug" env:"DEBUG" default:"false"`
	Version   string `long:"version" env:"VERSION" default:"unversioned"`
	Commit    string `long:"commit" env:"COMMIT"`
	BuildDate string `long:"build-date" env:"BUILD_DATE"`
	Name      string `long:"name" env:"NAME" default:"lazygit"`
}

// AppConfigurer interface allows individual app config structs to inherit Fields
// from AppConfig and still be used by lazygit.
type AppConfigurer interface {
	GetDebug() bool
	GetVersion() string
	GetCommit() string
	GetBuildDate() string
	GetName() string
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
