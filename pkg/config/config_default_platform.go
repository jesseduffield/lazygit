// +build !windows,!linux

package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	return OSConfig{
		EditCommand:     ``,
		OpenCommand:     "open {{filename}}",
		OpenLinkCommand: "open {{link}}",
	}
}
