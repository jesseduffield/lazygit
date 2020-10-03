package config

// OSConfig contains config on the level of the os
type OSConfig struct {
	// OpenCommand is the command for opening a file
	OpenCommand string `yaml:"openCommand,omitempty"`

	// OpenCommand is the command for opening a link
	OpenLinkCommand string `yaml:"openLinkCommand,omitempty"`
}
