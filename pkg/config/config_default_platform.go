// +build !windows,!linux

package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	return OSConfig{
		OpenCommand:     "open {{filename}}",
		OpenLinkCommand: "open {{link}}",
	}
}
